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

import {css, html, LitElement} from 'lit-element';
import {codePanelsStyles} from './styles';
import {basicSetup, EditorView} from 'codemirror';
import {Prec} from '@codemirror/state';
import {keymap} from '@codemirror/view';
import {indentWithTab, insertNewlineAndIndent} from '@codemirror/commands';
import {yaml} from '@codemirror/lang-yaml';
// forceLinting doesn't work reliably, using document change workaround instead
import {nothing} from 'lit';
import {repeat} from 'lit/directives/repeat.js';
import {
  StateField,
  StateEffect,
  RangeSet,
  EditorState,
  RangeSetBuilder,
  Compartment,
} from '@codemirror/state';
import {gutter, GutterMarker, Decoration, ViewPlugin} from '@codemirror/view';
import {ottlAutocompletion, ottlLinting} from '../ottl-completion.js';

export class PlaygroundConfigPanel extends LitElement {
  static properties = {
    examples: {type: Array},
    hideExamples: {type: Boolean, attribute: 'hide-examples'},
    config: {type: String},
    configDocsURL: {type: String, attribute: 'config-docs-url'},
    readOnly: {type: Boolean, attribute: 'read-only'},
    debuggerEnabled: {type: Boolean, attribute: 'debugger-enabled'},
    debuggingInfo: {type: Object, attribute: 'debugging-info'},

    _debuggingLineIndex: {state: true, type: Number},
    _debuggingLine: {state: true, type: Number},
    _editor: {state: true},
  };

  constructor() {
    super();
    this.hideExamples = false;
    this.examples = [];
    this.debuggerEnabled = true;
    this.configDocsURL = '';
    this._debuggingLineOffset = null;
    this._breakpointState = null;
    this._editorReadOnlyCompartment = new Compartment();
    this._editorBreakpointGutterCompartment = new Compartment();
    this.debuggingInfo = {
      debugging: false,
      lines: [],
      lineResultIndex: {},
    };
  }

  static get styles() {
    let styles = css`
      .debugger-controls {
        gap: 1px !important;
        display: inline-flex;
        align-items: center;
        justify-content: flex-end;
        padding: 5px 8px 5px 5px;
        border-left: gray 4px solid !important;
        border-bottom: #eee 1px solid;
      }

      .debugger-controls button {
        font-size: 13px;
        min-height: 25px;
        max-height: 25px;
        max-width: 32px;
        cursor: pointer;
        border: 2px;
        border-radius: 7px;
        display: flex;
        align-items: center;
      }

      .debugger-controls button:hover:enabled {
        background-color: #dedede;
      }
    `;

    return [...codePanelsStyles, styles];
  }

  get config() {
    return this._editor?.state.doc.toString() ?? '';
  }

  set config(value) {
    if (value === this.config) return;
    this.updateComplete.then(() => {
      this._editor?.dispatch({
        changes: {from: 0, to: this._editor.state.doc.length, insert: value},
      });
    });
  }

  firstUpdated() {
    this._initCodeEditor();
  }

  updated(changedProperties) {
    if (changedProperties.has('examples')) {
      // Reset the selected example
      this.shadowRoot.querySelector('#example-input').value = '';
    }

    if (changedProperties.has('debuggerEnabled')) {
      this._updateBreakpointGutterVisibility();
    }

    if (changedProperties.has('debuggingInfo')) {
      this._updateActiveDebugging();
    }

    if (changedProperties.has('_debuggingLine')) {
      this._refreshHighlightedDebuggingLine();
    }

    super.updated(changedProperties);
  }

  render() {
    return html`
      <div class="code-panel-parent" id="input-config-panel">
        <div class="code-panel-controls">
          <div class="header">
            <span>
              <strong>Configuration</strong>
              <sup>
                <a target="_blank" href="${this.configDocsURL}">
                  <small>YAML</small>
                </a>
              </sup>
            </span>
          </div>
          <div class="right" style="display: flex">
            ${this.hideExamples
              ? nothing
              : html`
                  <select
                    id="example-input"
                    @change="${this._handleExampleChanged}"
                    title="Select an example"
                    style="max-width:250px"
                    ?disabled="${this.debuggingInfo?.debugging === true}"
                  >
                    <option selected disabled value="">
                      Example ${'\u00A0'.repeat(45)}
                    </option>
                    ${this.examples &&
                    repeat(
                      this.examples,
                      (it) => it.name,
                      (it, idx) => {
                        return html`<option value="${idx}" title="${it.name}">
                          ${it.name}
                        </option>`;
                      }
                    )}
                  </select>
                `}
            <slot name="custom-components"></slot>
          </div>
        </div>
        ${this.debuggingInfo?.debugging
          ? html`
              <div class="debugger-controls">
                <button
                  @click="${this._stopDebuggingClick}"
                  title="Stop (&#8679;+F2)"
                >
                  <!-- prettier-ignore -->
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 -960 960 960"><path fill="darkred" d="M240-240v-480h480v480H240Z"/></svg>
                </button>
                <button
                  title="${this._hasNextDebuggingLine()
                    ? 'Resume'
                    : 'Rerun'} (&#8679;+F9)"
                  @click="${this._resumeDebuggingClick}"
                >
                  ${this._hasNextDebuggingLine()
                    ? html`
                        <!-- prettier-ignore -->
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 -960 960 960"><path fill="green" d="M240-240v-480h60v480h-60Zm174 0 385-240-385-240v480Z"/></svg>
                      `
                    : html`
                        <!-- prettier-ignore -->
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 -960 960 960"><path d="M451-122q-123-10-207-101t-84-216q0-77 35.5-145T295-695l43 43q-56 33-87 90.5T220-439q0 100 66 173t165 84v60Zm60 0v-60q100-12 165-84.5T741-439q0-109-75.5-184.5T481-699h-20l60 60-43 43-133-133 133-133 43 43-60 60h20q134 0 227 93.5T801-439q0 125-83.5 216T511-122Z"/></svg>
                      `}
                </button>
                <button
                  title="Step over (&#8679;+F8)"
                  @click="${this._nextDebugLineClick}"
                  ?disabled="${this._debuggingLineOffset ===
                  this.debuggingInfo?.lines?.length}"
                >
                  <!-- prettier-ignore -->
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 -960 960 960"><path d="M479.88-80Q434-80 402-112.12q-32-32.12-32-78T402.12-268q32.12-32 78-32T558-267.88q32 32.12 32 78T557.88-112q-32.12 32-78 32Zm.12-330L294-596l42-42 114 113v-354h60v354l113-113 43 42-186 186Z"/></svg>
                </button>
                <button
                  title="Step back (&#8679;+F7)"
                  @click="${this._previousDebugLineClick}"
                  ?disabled="${!this._debuggingLineOffset}"
                >
                  <!-- prettier-ignore -->
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 -960 960 960"><path d="M479.88-80Q434-80 402-112.12q-32-32.12-32-78T402.12-268q32.12-32 78-32T558-267.88q32 32.12 32 78T557.88-112q-32.12 32-78 32ZM450-410v-354L336-651l-42-42 186-186 186 186-43 42-113-113v354h-60Z"/></svg>
                </button>
              </div>
            `
          : nothing}
        <div class="code-editor-container">
          <div class="wrapper" id="config-input"></div>
        </div>
      </div>
    `;
  }

  hasBreakpoints() {
    if (
      !this._editor ||
      !this._editor.state.doc.length ||
      !this._breakpointState
    ) {
      return false;
    }

    let state = this._editor.state.field(this._breakpointState, false);
    if (!state) return false;

    let hasBreakpoints = false;
    for (let i = 1; i <= this._editor.state.doc.lines; i++) {
      let line = this._editor.state.doc.line(i);
      state.between(line.from, line.from, () => {
        hasBreakpoints = true;
      });
      if (hasBreakpoints) break;
    }
    return hasBreakpoints;
  }

  _updateBreakpointGutterVisibility() {
    this._editor?.dispatch({
      effects: this._editorBreakpointGutterCompartment.reconfigure(
        this.debuggerEnabled ? this.breakpointGutter : []
      ),
    });
  }

  _updateActiveDebugging() {
    this._debuggingLineOffset = this.debuggingInfo?.lines?.length;
    this._debuggingLine = null;

    let debugging = this.debuggingInfo?.debugging === true;
    if (debugging) {
      this._editor.dispatch({
        effects: this._editorReadOnlyCompartment.reconfigure(
          EditorState.readOnly.of(true)
        ),
      });
      this.updateComplete.then(() => {
        this._resumeDebuggingClick();
      });
    } else {
      if (this._editor?.state.doc.length > 0) {
        this._editor.dispatch({selection: {anchor: 1}, scrollIntoView: true});
      }
    }
  }

  _refreshHighlightedDebuggingLine() {
    if (this.debuggingInfo) {
      let anchor = null;
      if (this._debuggingLine === null || this._debuggingLine === undefined) {
        anchor = this._editor.state.selection?.main?.from || 0;
      } else if (
        this._debuggingLine > 0 &&
        this._editor.state.doc.lines >= this._debuggingLine
      ) {
        let line = this._editor.state.doc.line(this._debuggingLine);
        anchor = line.from;
      }
      if (anchor !== null) {
        this._editor.dispatch({
          selection: {anchor: anchor},
          scrollIntoView: true,
        });
      }
    }
  }

  _hasNextDebuggingLine() {
    return (
      this._debuggingLineOffset !== this.debuggingInfo?.lines?.length &&
      this._debuggingLine !== null
    );
  }

  _handleExampleChanged(event) {
    if (!event.target.value) return;
    let idx = parseInt(event.target.value);
    let example = this.examples[idx];
    if (!example) return;

    this.config = example.config;
    this.dispatchEvent(
      new CustomEvent('config-example-changed', {
        detail: {value: example},
        bubbles: true,
        composed: true,
        cancelable: true,
      })
    );
  }

  _resumeDebuggingClick() {
    let state = this._editor.state.field(this._breakpointState);
    let breakpoints = {};
    for (let i = 1; i <= this._editor.state.doc.lines; i++) {
      let line = this._editor.state.doc.line(i);
      state.between(line.from, line.from, () => {
        breakpoints[i] = true;
      });
    }

    if (Object.keys(breakpoints).length > 0) {
      for (let i = 0; i < this.debuggingInfo?.lines?.length; i++) {
        let line = this.debuggingInfo.lines[i];
        if (
          breakpoints[line] === true &&
          (this._debuggingLine == null || line > this._debuggingLine)
        ) {
          this._debuggingLine = line;
          this._debuggingLineOffset = i;
          this._notifyDebuggingLineChange(
            this.debuggingInfo.lines[i - 1] || -1
          );
          return;
        }
      }
    }

    this._resumeDebugging();
  }

  _resumeDebugging() {
    this._debuggingLineOffset = this.debuggingInfo?.lines?.length;
    this._debuggingLine = null;

    // Move to the end of the document
    let from = this._editor.state.doc.line(this._editor.state.doc.lines).from;
    this._editor.dispatch({
      selection: {anchor: from},
      scrollIntoView: true,
    });

    this._notifyDebuggingLineChange(
      this.debuggingInfo?.lines?.[this._debuggingLineOffset - 1]
    );

    this._stopDebuggingClick();
  }

  _nextDebugLineClick() {
    if (
      this._debuggingLineOffset === null ||
      this._debuggingLineOffset === undefined
    ) {
      this._debuggingLineOffset = 0;
      this._debuggingLine = this.debuggingInfo?.lines?.[0];
      return;
    }

    let nextIndex = this._debuggingLineOffset + 1;
    if (nextIndex > this.debuggingInfo?.lines?.length) {
      return;
    }

    this._debuggingLineOffset = nextIndex;
    if (nextIndex < this.debuggingInfo?.lines?.length) {
      this._debuggingLine = this.debuggingInfo?.lines?.[nextIndex];
    } else {
      this._debuggingLine = null;
      this._editor.dispatch({
        selection: {
          anchor: this._editor.state.doc.line(
            this.debuggingInfo?.lines?.[nextIndex - 1]
          ).from,
        },
        scrollIntoView: true,
      });
    }

    this._notifyDebuggingLineChange(
      this.debuggingInfo?.lines?.[this._debuggingLineOffset - 1]
    );
  }

  _previousDebugLineClick() {
    if (this._debuggingLineOffset === null || this._debuggingLineOffset <= 0) {
      return;
    }

    this._debuggingLineOffset--;
    this._debuggingLine =
      this.debuggingInfo?.lines?.[this._debuggingLineOffset];
    let beforePrevLine =
      this._debuggingLineOffset === 0
        ? -1
        : this.debuggingInfo?.lines?.[this._debuggingLineOffset - 1];

    this._notifyDebuggingLineChange(beforePrevLine);
  }

  _stopDebuggingClick() {
    this.debuggingInfo.debugging = false;
    this._debuggingLine = null;
    this._debuggingLineOffset = null;

    let readOnly = this.readOnly === true;
    if (this._editor?.state?.readOnly !== readOnly) {
      this._editor.dispatch({
        effects: this._editorReadOnlyCompartment.reconfigure(
          EditorState.readOnly.of(readOnly)
        ),
      });
    }

    this.dispatchEvent(
      new CustomEvent('debugging-stop-requested', {
        detail: {},
        bubbles: true,
        composed: true,
      })
    );
  }

  _notifyConfigChange(value) {
    this.dispatchEvent(
      new CustomEvent('config-changed', {
        detail: {value: value},
        bubbles: true,
        composed: true,
      })
    );
  }

  _notifyDebuggingLineChange(value) {
    this.dispatchEvent(
      new CustomEvent('debugging-line-changed', {
        detail: {value: value},
        bubbles: true,
        composed: true,
      })
    );
  }

  _initCodeEditor() {
    let me = this;

    const breakpointMarker = new (class extends GutterMarker {
      toDOM() {
        return document.createTextNode('⬤');
      }
    })();

    const breakpointEffect = StateEffect.define({
      map: (val, mapping) => ({pos: mapping.mapPos(val.pos), on: val.on}),
    });

    this._breakpointState = StateField.define({
      create() {
        return RangeSet.empty;
      },
      update(set, transaction) {
        set = set.map(transaction.changes);
        for (let e of transaction.effects) {
          if (e.is(breakpointEffect)) {
            if (e.value.on && me.debuggerEnabled === true)
              set = set.update({add: [breakpointMarker.range(e.value.pos)]});
            else set = set.update({filter: (from) => from !== e.value.pos});
          }
        }
        return set;
      },
    });

    const toggleBreakpoint = function (view, pos) {
      let breakpoints = view.state.field(me._breakpointState);
      let hasBreakpoint = false;
      breakpoints.between(pos, pos, () => {
        hasBreakpoint = true;
      });
      view.dispatch({
        effects: breakpointEffect.of({pos, on: !hasBreakpoint}),
      });
    };

    this.breakpointGutter = [
      this._breakpointState,
      gutter({
        class: 'cm-breakpoint-gutter',
        markers: (v) => v.state.field(me._breakpointState),
        initialSpacer: () => breakpointMarker,
        domEventHandlers: {
          mousedown(view, line) {
            toggleBreakpoint(view, line.from);
            return true;
          },
        },
      }),

      EditorView.baseTheme({
        '.cm-breakpoint-gutter': {
          cursor: 'pointer',
        },
        '.cm-breakpoint-gutter .cm-gutterElement': {
          color: 'red',
          paddingLeft: '3px',
          paddingRight: '3px',
          fontSize: '12px',
          cursor: 'default',
        },
      }),
    ];

    const debuggingTheme = EditorView.baseTheme({
      '&light .cm-debuggingLine': {
        backgroundColor: '#0a43c5',
        color: 'white!important',
      },
      '&light .cm-debuggingLine span': {
        color: 'white!important',
      },
    });

    const debuggingLineDeco = Decoration.line({
      attributes: {class: 'cm-debuggingLine'},
    });

    const debuggingDeco = function (view) {
      let builder = new RangeSetBuilder();
      if (
        me._debuggingLine != null &&
        view.state.doc.lines >= me._debuggingLine
      ) {
        let line = view.state.doc.line(me._debuggingLine);
        builder.add(line.from, line.from, debuggingLineDeco);
      }
      return builder.finish();
    };

    const showDebuggingLine = ViewPlugin.fromClass(
      class {
        constructor(view) {
          this.decorations = debuggingDeco(view);
        }
        update(update) {
          if (update.selectionSet)
            this.decorations = debuggingDeco(update.view);
        }
      },
      {
        decorations: (v) => v.decorations,
      }
    );

    const debuggingLinesExt = [debuggingTheme, showDebuggingLine];
    const debuggingShortcutsEnabled = () => {
      return (
        me.debuggerEnabled === true && me.debuggingInfo?.debugging === true
      );
    };

    let readOnly =
      this.readOnly === true || this.debuggingInfo?.debugging === true;
    this._editor = new EditorView({
      extensions: [
        basicSetup,
        Prec.highest(
          keymap.of([
            indentWithTab,
            {key: 'Enter', run: insertNewlineAndIndent, shift: () => true},
            {
              key: 'Shift-F2',
              run: () => {
                if (debuggingShortcutsEnabled()) this._stopDebuggingClick();
                return true;
              },
            },
            {
              key: 'Shift-F7',
              run: () => {
                if (debuggingShortcutsEnabled()) this._previousDebugLineClick();
                return true;
              },
            },
            {
              key: 'Shift-F8',
              run: () => {
                if (debuggingShortcutsEnabled()) this._nextDebugLineClick();
                return true;
              },
            },
            {
              key: 'Shift-F9',
              run: () => {
                if (debuggingShortcutsEnabled()) this._resumeDebuggingClick();
                return true;
              },
            },
          ])
        ),
        EditorView.lineWrapping,
        yaml(),
        this._editorBreakpointGutterCompartment.of(this.breakpointGutter),
        this._editorReadOnlyCompartment.of(EditorState.readOnly.of(readOnly)),
        debuggingLinesExt,
        ottlAutocompletion(),
        ottlLinting(),
        EditorView.updateListener.of((v) => {
          if (v.docChanged) {
            this._notifyConfigChange(this.config);
          }
        }),
      ],
      parent: this.shadowRoot.querySelector('#config-input'),
    });
  }

  /**
   * Force re-linting the editor content.
   * Call this when WASM version changes to re-validate with new validator.
   */
  forceLint() {
    if (this._editor) {
      // Workaround: Insert and remove a space to force document change
      // This triggers the linter to see a new document version
      // (forceLinting from @codemirror/lint doesn't work reliably)
      const docLength = this._editor.state.doc.length;
      this._editor.dispatch({
        changes: {from: docLength, insert: ' '},
      });
      this._editor.dispatch({
        changes: {from: docLength, to: docLength + 1, insert: ''},
      });
    }
  }
}

customElements.define('playground-config-panel', PlaygroundConfigPanel);
