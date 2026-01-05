/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. Elasticsearch B.V. licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import {autocompletion, startCompletion} from '@codemirror/autocomplete';
import {syntaxTree} from '@codemirror/language';
import {keymap, showTooltip, hoverTooltip, EditorView} from '@codemirror/view';
import {StateField} from '@codemirror/state';
import {linter} from '@codemirror/lint';

// Base URL for OTTL function documentation
const OTTL_DOCS_BASE =
  'https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl/ottlfuncs';

// Cache for OTTL metadata
let metadataCache = null;

/**
 * Clear the metadata cache. Call when WASM version changes.
 */
export function clearOTTLMetadataCache() {
  metadataCache = null;
}

/**
 * Get OTTL metadata from WASM, with caching.
 */
function getMetadata() {
  if (metadataCache) {
    return metadataCache;
  }

  // Check if WASM functions are available
  if (typeof window.getOTTLFunctions !== 'function') {
    return null;
  }

  try {
    metadataCache = {
      functions: window.getOTTLFunctions() || [],
      contexts: {},
    };

    // Load context-specific data for common contexts
    for (const ctx of ['log', 'span', 'metric', 'datapoint', 'resource', 'scope']) {
      metadataCache.contexts[ctx] = {
        paths: window.getContextPaths?.(ctx) || [],
        enums: window.getContextEnums?.(ctx) || [],
      };
    }
  } catch (e) {
    console.warn('Failed to load OTTL metadata:', e);
    return null;
  }

  return metadataCache;
}

/**
 * Detect the OTTL context from the YAML structure.
 * Looks for patterns like 'log_statements:', 'trace_statements:', 'context: log', etc.
 */
function detectOTTLContext(state, pos) {
  const doc = state.doc.toString();
  const beforePos = doc.slice(0, pos);

  // Check for explicit context declarations
  const contextMatch = beforePos.match(/context:\s*(\w+)\s*$/m);
  if (contextMatch) {
    const ctx = contextMatch[1].toLowerCase();
    if (ctx === 'span' || ctx === 'spanevent') return 'span';
    if (ctx === 'datapoint') return 'datapoint';
    return ctx;
  }

  // Check for statement type prefixes
  if (beforePos.includes('log_statements:') || beforePos.includes('log:')) {
    return 'log';
  }
  if (beforePos.includes('trace_statements:') || beforePos.includes('span:')) {
    return 'span';
  }
  if (beforePos.includes('metric_statements:') || beforePos.includes('metric:')) {
    return 'metric';
  }
  if (beforePos.includes('profile_statements:')) {
    return 'profile';
  }

  // Default to log context
  return 'log';
}

/**
 * Check if the cursor is inside a statements array in YAML.
 */
function isInsideStatementsArray(state, pos) {
  const tree = syntaxTree(state);
  let node = tree.resolveInner(pos, -1);

  // Walk up the tree to check if we're in a relevant context
  while (node) {
    // Check if we're in a sequence item (list item in YAML)
    if (node.name === 'BlockSequence' || node.name === 'FlowSequence') {
      // Check the parent context - are we under a 'statements' key?
      const lineText = state.doc.lineAt(node.from).text;
      if (lineText.includes('statements:')) {
        return true;
      }

      // Check lines before to see if we're under statements
      const textBefore = state.doc.sliceString(0, node.from);
      const lines = textBefore.split('\n').slice(-10);
      for (const line of lines) {
        if (line.match(/^\s*statements:\s*$/)) {
          return true;
        }
        if (line.match(/^\s*-\s+context:/)) {
          // We're in a transform processor config
          return true;
        }
      }
    }
    node = node.parent;
  }

  // Fallback: check if line starts with '- ' (list item) and nearby lines suggest statements
  const line = state.doc.lineAt(pos);
  const lineText = line.text;
  if (lineText.match(/^\s*-\s/)) {
    const textBefore = state.doc.sliceString(Math.max(0, pos - 500), pos);
    if (textBefore.includes('statements:')) {
      return true;
    }
  }

  return false;
}

/**
 * Generate completion items for OTTL functions.
 */
function functionCompletions(functions) {
  return functions.map((f) => {
    const params = f.parameters
      .map((p) => {
        const optional = p.optional ? '?' : '';
        return `${p.name}${optional}`;
      })
      .join(', ');

    // Create snippet for required parameters
    const requiredParams = f.parameters.filter((p) => !p.optional);
    const snippet =
      requiredParams.length > 0
        ? `${f.name}(${requiredParams.map((p, i) => `\${${i + 1}:${p.name}}`).join(', ')})`
        : `${f.name}()`;

    return {
      label: f.name,
      type: f.isEditor ? 'function' : 'method',
      detail: `(${params})`,
      info: f.isEditor ? 'Editor - modifies data' : 'Converter - returns value',
      apply: snippet,
      boost: f.isEditor ? 2 : 1, // Prioritize editors at statement start
    };
  });
}

/**
 * Generate completion items for context paths.
 */
function pathCompletions(paths) {
  return paths.map((p) => ({
    label: p.path,
    type: 'property',
    detail: p.type,
    info: p.description || (p.supportsKeys ? 'Supports key access: ["key"]' : ''),
  }));
}

/**
 * Generate completion items for enums.
 */
function enumCompletions(enums) {
  return enums.map((e) => ({
    label: e.name,
    type: 'constant',
    detail: `= ${e.value}`,
  }));
}

/**
 * Main completion source for OTTL.
 */
function ottlCompletionSource(context) {
  // Check if we're inside a statements array
  if (!isInsideStatementsArray(context.state, context.pos)) {
    return null;
  }

  const metadata = getMetadata();
  if (!metadata) {
    return null;
  }

  // Detect the OTTL context (log, span, metric, etc.)
  const ottlContext = detectOTTLContext(context.state, context.pos);
  const contextData = metadata.contexts[ottlContext] || metadata.contexts['log'];

  // Get the word being typed
  const word = context.matchBefore(/[\w.]*/);
  const line = context.state.doc.lineAt(context.pos);
  const textBefore = line.text.slice(0, context.pos - line.from);

  // Determine what kind of completions to provide
  let completions = [];

  // Check if we're typing a path continuation (after a dot)
  if (textBefore.match(/\.\s*$/)) {
    // Path continuation - filter paths by prefix
    const match = textBefore.match(/([\w.]+)\.\s*$/);
    if (match) {
      const prefix = match[1] + '.';
      completions = pathCompletions(
        contextData.paths.filter((p) => p.path.startsWith(prefix))
      );
    }
  }
  // At the start of a statement (after '- ' or '- "')
  else if (textBefore.match(/^\s*-\s*['"]?\s*$/)) {
    // Suggest functions (primarily editors for statements)
    completions = functionCompletions(metadata.functions);
  }
  // Inside function arguments
  else if (textBefore.includes('(')) {
    const openParens = (textBefore.match(/\(/g) || []).length;
    const closeParens = (textBefore.match(/\)/g) || []).length;

    if (openParens > closeParens) {
      // Inside function arguments - suggest paths, enums, and converters
      completions = [
        ...pathCompletions(contextData.paths),
        ...enumCompletions(contextData.enums || []),
        ...functionCompletions(metadata.functions.filter((f) => !f.isEditor)),
      ];
    }
  }
  // Default - show all completions
  else {
    completions = [
      ...functionCompletions(metadata.functions),
      ...pathCompletions(contextData.paths),
      ...enumCompletions(contextData.enums || []),
    ];
  }

  if (completions.length === 0) {
    return null;
  }

  return {
    from: word ? word.from : context.pos,
    options: completions,
    validFor: /^[\w.]*$/,
  };
}

/**
 * Parse the current function call context from the text before cursor.
 * Returns {functionName, argIndex} or null if not inside a function call.
 */
function parseFunctionContext(textBefore) {
  // Find the last unmatched opening parenthesis
  let depth = 0;
  let funcStart = -1;

  for (let i = textBefore.length - 1; i >= 0; i--) {
    const char = textBefore[i];
    if (char === ')') {
      depth++;
    } else if (char === '(') {
      if (depth === 0) {
        funcStart = i;
        break;
      }
      depth--;
    }
  }

  if (funcStart === -1) {
    return null;
  }

  // Extract function name (word before the opening paren)
  const beforeParen = textBefore.slice(0, funcStart);
  const funcMatch = beforeParen.match(/(\w+)\s*$/);
  if (!funcMatch) {
    return null;
  }

  const functionName = funcMatch[1];

  // Count arguments (by counting commas at depth 0)
  const insideParens = textBefore.slice(funcStart + 1);
  let argIndex = 0;
  depth = 0;

  for (const char of insideParens) {
    if (char === '(' || char === '[') {
      depth++;
    } else if (char === ')' || char === ']') {
      depth--;
    } else if (char === ',' && depth === 0) {
      argIndex++;
    }
  }

  return {functionName, argIndex};
}

/**
 * Get function info by name from metadata.
 */
function getFunctionByName(name) {
  const metadata = getMetadata();
  if (!metadata) return null;

  // Case-insensitive search
  const lowerName = name.toLowerCase();
  return metadata.functions.find(
    (f) => f.name.toLowerCase() === lowerName
  );
}

/**
 * Create the signature help tooltip content.
 */
function createSignatureTooltip(func, argIndex) {
  const container = document.createElement('div');
  container.className = 'cm-signature-help';
  container.style.cssText = `
    padding: 4px 8px;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: 4px;
    font-family: monospace;
    font-size: 13px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.3);
    max-width: 500px;
  `;

  // Build signature with highlighted current parameter
  let html = `<span style="color: #dcdcaa">${func.name}</span>(`;

  func.parameters.forEach((param, i) => {
    if (i > 0) html += ', ';

    const isActive = i === argIndex;
    const style = isActive
      ? 'color: #9cdcfe; font-weight: bold; text-decoration: underline;'
      : 'color: #9cdcfe;';

    const optional = param.optional ? '?' : '';
    html += `<span style="${style}">${param.name}${optional}</span>`;
  });

  html += ')';

  // Add parameter description if available
  if (func.parameters[argIndex]) {
    const param = func.parameters[argIndex];
    html += `<div style="margin-top: 4px; color: #808080; font-size: 12px;">`;
    html += `<strong>${param.name}</strong>: ${param.kind}`;
    if (param.optional) html += ' (optional)';
    html += '</div>';
  }

  container.innerHTML = html;
  return container;
}

/**
 * StateField to track signature help tooltip.
 */
const signatureHelpState = StateField.define({
  create() {
    return null;
  },

  update(value, tr) {
    // Only update on document changes or selection changes
    if (!tr.docChanged && !tr.selection) {
      return value;
    }

    const pos = tr.state.selection.main.head;
    const line = tr.state.doc.lineAt(pos);
    const textBefore = line.text.slice(0, pos - line.from);

    // Check if we're in a statements context first
    if (!isInsideStatementsArray(tr.state, pos)) {
      return null;
    }

    const context = parseFunctionContext(textBefore);
    if (!context) {
      return null;
    }

    const func = getFunctionByName(context.functionName);
    if (!func || func.parameters.length === 0) {
      return null;
    }

    return {
      pos: pos,
      func: func,
      argIndex: Math.min(context.argIndex, func.parameters.length - 1),
    };
  },

  provide(field) {
    return showTooltip.compute([field], (state) => {
      const data = state.field(field);
      if (!data) return null;

      return {
        pos: data.pos,
        above: true,
        strictSide: true,
        arrow: false,
        create() {
          return {dom: createSignatureTooltip(data.func, data.argIndex)};
        },
      };
    });
  },
});

// Store for validation context (data type, executor, and payload)
let validationContext = {
  dataType: 'logs',
  executor: 'transform',
  payload: '{}',
};

/**
 * Set the validation context for the linter.
 * Call this when the user changes data type, executor, or payload.
 */
export function setValidationContext(dataType, executor, payload = '{}') {
  validationContext = {dataType, executor, payload};
}

/**
 * OTTL linter that validates statements via WASM.
 */
function ottlLinter(view) {
  const diagnostics = [];

  // Check if validation function is available
  if (typeof window.validateStatements !== 'function') {
    return diagnostics;
  }

  const doc = view.state.doc.toString();

  // Only validate if there's content
  if (!doc.trim()) {
    return diagnostics;
  }

  try {
    const results = window.validateStatements(
      doc,
      validationContext.dataType,
      validationContext.payload,
      validationContext.executor
    );

    if (!results || !Array.isArray(results)) {
      return diagnostics;
    }

    for (const diag of results) {
      // Convert line/column to position
      const line = Math.max(1, diag.line || 1);
      const lineInfo = view.state.doc.line(Math.min(line, view.state.doc.lines));
      const col = Math.max(0, (diag.column || 1) - 1);
      const from = lineInfo.from + Math.min(col, lineInfo.length);

      const endLine = Math.max(1, diag.endLine || line);
      const endLineInfo = view.state.doc.line(
        Math.min(endLine, view.state.doc.lines)
      );
      const endCol = Math.max(0, (diag.endColumn || col + 1) - 1);
      const to = endLineInfo.from + Math.min(endCol, endLineInfo.length);

      diagnostics.push({
        from: from,
        to: Math.max(from, to),
        severity: diag.severity === 'warning' ? 'warning' : 'error',
        message: diag.message,
      });
    }
  } catch (e) {
    console.warn('OTTL validation error:', e);
  }

  return diagnostics;
}

/**
 * Create documentation content for a function.
 */
function createFunctionDoc(func) {
  const container = document.createElement('div');
  container.className = 'cm-hover-doc';
  container.style.cssText = `
    padding: 8px 12px;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: 4px;
    font-family: monospace;
    font-size: 13px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.3);
    max-width: 500px;
  `;

  // Function signature
  const params = func.parameters
    .map((p) => {
      const optional = p.optional ? '?' : '';
      return `${p.name}${optional}: ${p.kind}`;
    })
    .join(', ');

  let html = `<div style="color: #dcdcaa; font-weight: bold; margin-bottom: 4px;">${func.name}(${params})</div>`;

  // Function type
  html += `<div style="color: #808080; font-size: 12px; margin-bottom: 8px;">`;
  html += func.isEditor ? 'Editor - modifies telemetry data' : 'Converter - returns a value';
  html += '</div>';

  // Parameter details
  if (func.parameters.length > 0) {
    html += '<div style="border-top: 1px solid #444; padding-top: 6px; margin-top: 4px;">';
    html += '<div style="color: #569cd6; font-size: 11px; margin-bottom: 4px;">Parameters:</div>';
    for (const p of func.parameters) {
      html += `<div style="margin-left: 8px; font-size: 12px;">`;
      html += `<span style="color: #9cdcfe;">${p.name}</span>`;
      html += `<span style="color: #808080;">: ${p.kind}</span>`;
      if (p.optional) html += '<span style="color: #6a9955;"> (optional)</span>';
      if (p.isSlice) html += '<span style="color: #ce9178;"> []</span>';
      html += '</div>';
    }
    html += '</div>';
  }

  container.innerHTML = html;
  return container;
}

/**
 * Create documentation content for a path.
 */
function createPathDoc(path) {
  const container = document.createElement('div');
  container.className = 'cm-hover-doc';
  container.style.cssText = `
    padding: 8px 12px;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: 4px;
    font-family: monospace;
    font-size: 13px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.3);
    max-width: 400px;
  `;

  let html = `<div style="color: #4ec9b0; font-weight: bold; margin-bottom: 4px;">${path.path}</div>`;
  html += `<div style="color: #808080; font-size: 12px;">Type: ${path.type}</div>`;

  if (path.description) {
    html += `<div style="margin-top: 6px; color: #d4d4d4;">${path.description}</div>`;
  }

  if (path.supportsKeys) {
    html += `<div style="margin-top: 4px; color: #6a9955; font-size: 11px;">Supports key access: ["key"]</div>`;
  }

  container.innerHTML = html;
  return container;
}

/**
 * Create documentation content for an enum.
 */
function createEnumDoc(enumInfo) {
  const container = document.createElement('div');
  container.className = 'cm-hover-doc';
  container.style.cssText = `
    padding: 8px 12px;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: 4px;
    font-family: monospace;
    font-size: 13px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.3);
    max-width: 300px;
  `;

  let html = `<div style="color: #b5cea8; font-weight: bold;">${enumInfo.name}</div>`;
  html += `<div style="color: #808080; font-size: 12px;">Value: ${enumInfo.value}</div>`;

  container.innerHTML = html;
  return container;
}

/**
 * Get documentation URL for a function.
 */
function getFunctionDocsUrl(funcName) {
  // OTTL functions are documented in ottlfuncs README with anchors
  // GitHub anchors use lowercase, underscores stay as underscores
  const anchor = funcName.toLowerCase();
  return `${OTTL_DOCS_BASE}#${anchor}`;
}

/**
 * Get word and metadata at a given position.
 */
function getWordAtPosition(view, pos) {
  if (!isInsideStatementsArray(view.state, pos)) {
    return null;
  }

  const metadata = getMetadata();
  if (!metadata) {
    return null;
  }

  const line = view.state.doc.lineAt(pos);
  const text = line.text;
  const linePos = pos - line.from;

  // Find word boundaries
  let start = linePos;
  let end = linePos;
  while (start > 0 && /[\w.]/.test(text[start - 1])) start--;
  while (end < text.length && /[\w.]/.test(text[end])) end++;

  const word = text.slice(start, end);
  if (!word) return null;

  // Check if it's a function
  const func = metadata.functions.find(
    (f) => f.name.toLowerCase() === word.toLowerCase()
  );
  if (func) {
    return {
      type: 'function',
      name: func.name,
      data: func,
      from: line.from + start,
      to: line.from + end,
    };
  }

  // Detect context for paths
  const ottlContext = detectOTTLContext(view.state, pos);
  const contextData = metadata.contexts[ottlContext] || metadata.contexts['log'];

  // Check if it's a path
  const path = contextData.paths.find((p) => p.path === word);
  if (path) {
    return {
      type: 'path',
      name: path.path,
      data: path,
      from: line.from + start,
      to: line.from + end,
    };
  }

  return null;
}

/**
 * Handle go-to-definition (Cmd/Ctrl+Click).
 */
function handleGoToDefinition(view, pos) {
  const wordInfo = getWordAtPosition(view, pos);
  if (!wordInfo) return false;

  if (wordInfo.type === 'function') {
    // Open function documentation
    const url = getFunctionDocsUrl(wordInfo.name);
    window.open(url, '_blank');
    return true;
  }

  // For paths, we could open context documentation
  // Currently just return false as paths don't have direct doc links
  return false;
}

/**
 * Extension for go-to-definition with Cmd/Ctrl+Click.
 */
const goToDefinitionExtension = EditorView.domEventHandlers({
  click(event, view) {
    // Check for Cmd (Mac) or Ctrl (Windows/Linux)
    if (!event.metaKey && !event.ctrlKey) {
      return false;
    }

    const pos = view.posAtCoords({x: event.clientX, y: event.clientY});
    if (pos === null) return false;

    if (handleGoToDefinition(view, pos)) {
      event.preventDefault();
      return true;
    }

    return false;
  },
});

/**
 * Add underline styling when hovering with Cmd/Ctrl pressed.
 */
const goToDefinitionHoverStyle = EditorView.domEventHandlers({
  mousemove(event, view) {
    // We add a class to indicate clickable state
    const editorDom = view.dom;
    if (event.metaKey || event.ctrlKey) {
      const pos = view.posAtCoords({x: event.clientX, y: event.clientY});
      if (pos !== null) {
        const wordInfo = getWordAtPosition(view, pos);
        if (wordInfo && wordInfo.type === 'function') {
          editorDom.style.cursor = 'pointer';
          return;
        }
      }
    }
    editorDom.style.cursor = '';
  },
  keydown(event) {
    if (event.key === 'Meta' || event.key === 'Control') {
      // Could add visual feedback here
    }
  },
  keyup(event, view) {
    if (event.key === 'Meta' || event.key === 'Control') {
      view.dom.style.cursor = '';
    }
  },
});

/**
 * Hover tooltip provider for OTTL elements.
 */
const ottlHoverTooltip = hoverTooltip((view, pos) => {
  // Check if we're in statements context
  if (!isInsideStatementsArray(view.state, pos)) {
    return null;
  }

  const metadata = getMetadata();
  if (!metadata) {
    return null;
  }

  // Get the word at position
  const line = view.state.doc.lineAt(pos);
  const text = line.text;
  const linePos = pos - line.from;

  // Find word boundaries
  let start = linePos;
  let end = linePos;
  while (start > 0 && /[\w.]/.test(text[start - 1])) start--;
  while (end < text.length && /[\w.]/.test(text[end])) end++;

  const word = text.slice(start, end);
  if (!word) return null;

  // Detect context for paths/enums
  const ottlContext = detectOTTLContext(view.state, pos);
  const contextData = metadata.contexts[ottlContext] || metadata.contexts['log'];

  // Check if it's a function
  const func = metadata.functions.find(
    (f) => f.name.toLowerCase() === word.toLowerCase()
  );
  if (func) {
    return {
      pos: line.from + start,
      end: line.from + end,
      above: true,
      create() {
        return {dom: createFunctionDoc(func)};
      },
    };
  }

  // Check if it's a path
  const path = contextData.paths.find((p) => p.path === word);
  if (path) {
    return {
      pos: line.from + start,
      end: line.from + end,
      above: true,
      create() {
        return {dom: createPathDoc(path)};
      },
    };
  }

  // Check if it's an enum
  const enumInfo = (contextData.enums || []).find((e) => e.name === word);
  if (enumInfo) {
    return {
      pos: line.from + start,
      end: line.from + end,
      above: true,
      create() {
        return {dom: createEnumDoc(enumInfo)};
      },
    };
  }

  return null;
});

/**
 * Create the OTTL autocompletion extension for CodeMirror.
 */
export function ottlAutocompletion() {
  return [
    autocompletion({
      override: [ottlCompletionSource],
      activateOnTyping: true,
      selectOnOpen: false,
      closeOnBlur: true,
      maxRenderedOptions: 50,
    }),
    // Signature help tooltip
    signatureHelpState,
    // Hover documentation
    ottlHoverTooltip,
    // Go-to-definition (Cmd/Ctrl+Click)
    goToDefinitionExtension,
    goToDefinitionHoverStyle,
    // Add macOS-friendly keybinding (Cmd+. or Escape then Tab)
    keymap.of([
      {key: 'Mod-.', run: startCompletion}, // Cmd+. on macOS, Ctrl+. on Windows/Linux
      {key: 'Alt-Space', run: startCompletion}, // Alt+Space as fallback
    ]),
  ];
}

/**
 * Create the OTTL linter extension for CodeMirror.
 * Call this separately so validation can be enabled/disabled.
 */
export function ottlLinting() {
  return linter(ottlLinter, {
    delay: 500, // Debounce validation by 500ms
  });
}
