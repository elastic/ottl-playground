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

import {Decoration, EditorView} from '@codemirror/view';
import {StateField, StateEffect} from '@codemirror/state';
import {syntaxTree} from '@codemirror/language';

const ottlClickableHoverRange = StateEffect.define({map: (v) => v});

const ottlClickableHoverField = StateField.define({
  create() {
    return Decoration.none;
  },
  update(set, tr) {
    set = set.map(tr.changes);
    for (const e of tr.effects) {
      if (e.is(ottlClickableHoverRange)) {
        if (e.value == null) {
          set = Decoration.none;
        } else {
          set = Decoration.set([
            Decoration.mark({
              class: 'cm-ottl-clickable-hover',
              inclusive: false,
            }).range(e.value.from, e.value.to),
          ]);
        }
      }
    }
    return set;
  },
  provide: (f) => EditorView.decorations.from(f),
});

function findClickableTokenAt(state, pos) {
  const tree = syntaxTree(state);
  let cur = tree.resolveInner(pos, 0);
  while (cur) {
    if (cur.name === 'Editor') {
      const first = cur.getChild('Lowercase');
      if (first) return {from: first.from, to: first.to, name: cur.name};
      return {from: cur.from, to: cur.to, name: cur.name};
    }
    if (cur.name === 'Converter') {
      const first = cur.getChild('Uppercase');
      if (first) return {from: first.from, to: first.to, name: cur.name};
      return {from: cur.from, to: cur.to, name: cur.name};
    }
    if (cur.name === 'Path') {
      const first = cur.getChild('Lowercase');
      const segments = cur.getChild('PathTail');
      // only context names should be clickable, not full paths.
      if (first && segments && segments.getChild('Lowercase')) {
        return {from: first.from, to: first.to, name: cur.name};
      }
    }
    cur = cur.parent;
  }
  return null;
}

export function ottlClickableHoverExtension(clickHandler = () => {}) {
  return [
    ottlClickableHoverField,
    EditorView.domEventHandlers({
      mousemove(ev, view) {
        const showUnderline = ev.ctrlKey || ev.metaKey;
        if (!showUnderline) {
          view.dispatch({effects: ottlClickableHoverRange.of(null)});
          return;
        }
        const pos = view.posAtCoords({x: ev.clientX, y: ev.clientY});
        if (pos == null) {
          view.dispatch({effects: ottlClickableHoverRange.of(null)});
          return;
        }
        const range = findClickableTokenAt(view.state, pos);
        view.dispatch({
          effects: ottlClickableHoverRange.of(
            range ? {from: range.from, to: range.to} : null
          ),
        });
      },
      mousedown(ev, view) {
        if (!ev.ctrlKey && !ev.metaKey) {
          return;
        }

        const pos = view.posAtCoords({x: ev.clientX, y: ev.clientY});
        if (pos == null) {
          return;
        }
        const range = findClickableTokenAt(view.state, pos);
        if (!range) {
          return;
        }

        ev.preventDefault();
        ev.stopImmediatePropagation();

        const text = view.state.sliceDoc(range.from, range.to);
        clickHandler({
          name: range.name,
          from: range.from,
          to: range.to,
          text,
        });
      },
      mouseleave(ev, view) {
        if (!view.dom.contains(ev.relatedTarget)) {
          view.dispatch({effects: ottlClickableHoverRange.of(null)});
        }
      },
      keyup(ev, view) {
        if (ev.key === 'Control' || ev.key === 'Meta') {
          view.dispatch({effects: ottlClickableHoverRange.of(null)});
        }
      },
    }),
    EditorView.theme({
      '.cm-ottl-clickable-hover': {
        textDecoration: 'underline',
        cursor: 'pointer',
      },
    }),
  ];
}
