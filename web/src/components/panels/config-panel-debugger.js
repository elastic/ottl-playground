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

import {
  RangeSet,
  RangeSetBuilder,
  StateEffect,
  StateField,
} from '@codemirror/state';
import {
  Decoration,
  EditorView,
  gutter,
  GutterMarker,
  ViewPlugin,
} from '@codemirror/view';

const breakpointIcon = '⬤';

const breakpointMarker = new (class extends GutterMarker {
  toDOM() {
    return document.createTextNode(breakpointIcon);
  }
})();

const hoverBreakpointMarker = new (class extends GutterMarker {
  toDOM() {
    const span = document.createElement('span');
    span.className = 'cm-breakpoint-gutter-hover';
    span.textContent = breakpointIcon;
    return span;
  }
})();

const breakpointEffect = StateEffect.define({
  map: (val, mapping) => ({pos: mapping.mapPos(val.pos), on: val.on}),
});

const hoveredLineEffect = StateEffect.define({
  map: (val, mapping) => (val != null ? mapping.mapPos(val) : null),
});

/**
 * Creates the config panel debugger extension: breakpoint gutter, current-line
 * highlight, and keyboard shortcuts. The panel passes an API object so the
 * extension can read state and trigger actions without depending on the panel.
 *
 * @param {{
 *   getDebuggerEnabled: () => boolean,
 *   getDebuggingLine: () => number | null,
 *   isDebugging: () => boolean,
 *   onStop: () => void,
 *   onResume: () => void,
 *   onNextLine: () => void,
 *   onPreviousLine: () => void,
 * }} api
 * @returns {{ breakpointState: StateField, breakpointGutter: unknown[], debuggingLineExt: unknown[], debuggerKeymap: unknown[] }}
 */
export function configPanelDebuggerExtension(api) {
  const breakpointState = StateField.define({
    create() {
      return RangeSet.empty;
    },
    update(set, transaction) {
      set = set.map(transaction.changes);
      for (const e of transaction.effects) {
        if (e.is(breakpointEffect)) {
          if (e.value.on && api.getDebuggerEnabled() === true) {
            set = set.update({add: [breakpointMarker.range(e.value.pos)]});
          } else {
            set = set.update({filter: (from) => from !== e.value.pos});
          }
        }
      }
      return set;
    },
  });

  const hoveredLineField = StateField.define({
    create() {
      return null;
    },
    update(value, transaction) {
      for (const e of transaction.effects) {
        if (e.is(hoveredLineEffect)) {
          value = e.value;
        }
      }
      return value;
    },
  });

  const toggleBreakpoint = function (view, pos) {
    const breakpoints = view.state.field(breakpointState);
    let hasBreakpoint = false;
    breakpoints.between(pos, pos, () => {
      hasBreakpoint = true;
    });
    view.dispatch({
      effects: breakpointEffect.of({pos, on: !hasBreakpoint}),
    });
  };

  const gutterMarkers = (view) => {
    const breakpoints = view.state.field(breakpointState);
    const hoveredFrom = view.state.field(hoveredLineField);
    if (hoveredFrom == null) {
      return breakpoints;
    }
    let hasBreakpointAtHover = false;
    breakpoints.between(hoveredFrom, hoveredFrom, () => {
      hasBreakpointAtHover = true;
    });
    if (hasBreakpointAtHover) {
      return breakpoints;
    }
    return breakpoints.update({
      add: [hoverBreakpointMarker.range(hoveredFrom)],
    });
  };

  const breakpointGutterHover = ViewPlugin.fromClass(
    class {
      constructor(view) {
        this.view = view;
        this.gutterEl = null;
        this.clear = () => {
          if (this.view.state.field(hoveredLineField) != null) {
            this.view.dispatch({effects: hoveredLineEffect.of(null)});
          }
        };
        this.onMousemove = (e) => {
          const pos = this.view.posAtCoords({x: e.clientX, y: e.clientY});
          if (pos == null) {
            return;
          }
          const line = this.view.state.doc.lineAt(pos);
          const current = this.view.state.field(hoveredLineField);
          if (current !== line.from) {
            this.view.dispatch({effects: hoveredLineEffect.of(line.from)});
          }
        };
      }
      update(update) {
        this.view = update.view;
        const el = update.view.dom.querySelector('.cm-breakpoint-gutter');
        if (el === this.gutterEl) {
          return;
        }
        if (this.gutterEl) {
          this.gutterEl.removeEventListener('mouseleave', this.clear);
          this.gutterEl.removeEventListener('mousemove', this.onMousemove);
        }
        this.gutterEl = el;
        if (el) {
          el.addEventListener('mouseleave', this.clear);
          el.addEventListener('mousemove', this.onMousemove);
        }
      }
      destroy() {
        if (this.gutterEl) {
          this.gutterEl.removeEventListener('mouseleave', this.clear);
          this.gutterEl.removeEventListener('mousemove', this.onMousemove);
          this.gutterEl = null;
        }
      }
    }
  );

  const breakpointGutter = [
    breakpointState,
    hoveredLineField,
    gutter({
      class: 'cm-breakpoint-gutter',
      markers: gutterMarkers,
      initialSpacer: () => breakpointMarker,
      domEventHandlers: {
        mousedown(view, line) {
          toggleBreakpoint(view, line.from);
          return true;
        },
        mouseover(view, line) {
          const current = view.state.field(hoveredLineField);
          if (current !== line.from) {
            view.dispatch({effects: hoveredLineEffect.of(line.from)});
          }
          return false;
        },
      },
    }),
    breakpointGutterHover,
    EditorView.baseTheme({
      '.cm-breakpoint-gutter': {
        cursor: 'pointer',
      },
      '.cm-breakpoint-gutter .cm-gutterElement': {
        color: 'red',
        paddingLeft: '3px',
        paddingRight: '3px',
        fontSize: '12px',
        cursor: 'pointer',
      },
      '.cm-breakpoint-gutter-hover': {
        color: 'red',
        opacity: '0.4',
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
    const builder = new RangeSetBuilder();
    const currentLine = api.getDebuggingLine();
    if (currentLine != null && view.state.doc.lines >= currentLine) {
      const line = view.state.doc.line(currentLine);
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
        if (update.selectionSet) {
          this.decorations = debuggingDeco(update.view);
        }
      }
    },
    {
      decorations: (v) => v.decorations,
    }
  );

  const debuggingLineExt = [debuggingTheme, showDebuggingLine];

  const debuggingShortcutsEnabled = () =>
    api.getDebuggerEnabled() === true && api.isDebugging() === true;

  const debuggerKeymap = [
    {
      key: 'Shift-F2',
      run: () => {
        if (debuggingShortcutsEnabled()) api.onStop();
        return true;
      },
    },
    {
      key: 'Shift-F7',
      run: () => {
        if (debuggingShortcutsEnabled()) api.onPreviousLine();
        return true;
      },
    },
    {
      key: 'Shift-F8',
      run: () => {
        if (debuggingShortcutsEnabled()) api.onNextLine();
        return true;
      },
    },
    {
      key: 'Shift-F9',
      run: () => {
        if (debuggingShortcutsEnabled()) api.onResume();
        return true;
      },
    },
  ];

  return {
    breakpointState,
    breakpointGutter,
    debuggingLineExt,
    debuggerKeymap,
  };
}
