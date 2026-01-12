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

import {html, LitElement} from 'lit-element';
import {codePanelsStyles} from './styles';
import {basicSetup, EditorView} from 'codemirror';
import {keymap} from '@codemirror/view';
import {indentWithTab, insertNewlineAndIndent} from '@codemirror/commands';
import {linter, lintGutter} from '@codemirror/lint';
import {json, jsonParseLinter} from '@codemirror/lang-json';
import {nothing} from 'lit';
import {Prec, Compartment, EditorState} from '@codemirror/state';
import {repeat} from 'lit/directives/repeat.js';

export class PlaygroundPayloadPanel extends LitElement {
  static properties = {
    payload: {type: String},
    examples: {type: Array},
    hideExamples: {type: Boolean, attribute: 'hide-examples'},
    readOnly: {type: Boolean, attribute: 'read-only'},
    _editor: {type: Object, state: true},
  };

  constructor() {
    super();
    this.payload = '{}';
    this.examples = [];
    this.hideExamples = false;
    this._editorReadOnlyCompartment = new Compartment();
  }

  static get styles() {
    return codePanelsStyles;
  }

  get payload() {
    return this._editor?.state.doc.toString() ?? '';
  }

  set payload(value) {
    if (value === this.payload) return;
    this.updateComplete.then(() => {
      this._editor?.dispatch({
        changes: {from: 0, to: this._editor.state.doc.length, insert: value},
      });
    });
  }

  clearSelectedExample() {
    this.shadowRoot.querySelector('#example-select').value = '';
  }

  firstUpdated() {
    this._initCodeEditor();
  }

  updated(changedProperties) {
    if (changedProperties.has('examples')) {
      this.clearSelectedExample();
    }

    if (changedProperties.has('readOnly')) {
      this._editor?.dispatch({
        effects: this._editorReadOnlyCompartment.reconfigure(
          EditorState.readOnly.of(this.readOnly || false)
        ),
      });
    }

    super.updated(changedProperties);
  }

  render() {
    return html`
      <div class="code-panel-parent" id="input-sample-panel">
        <div class="code-panel-controls">
          <div class="header">
            <span
              ><strong>OTLP payload</strong>
              <sup
                ><small
                  ><a
                    target="_blank"
                    href="https://opentelemetry.io/docs/specs/otlp/#json-protobuf-encoding"
                    >JSON</a
                  ></small
                ></sup
              ></span
            >
          </div>
          <div class="right" style="display: flex">
            ${this.hideExamples
              ? nothing
              : html`
                                    <select
                                            id="example-select"
                                            @change="${this._handleExampleChanged}"
                                            title="Payload example"
                                            ?disabled="${this.readOnly === true}"
                                    >
                                        <option selected disabled value="">Example</option>
                                      ${
                                        this.examples &&
                                        repeat(
                                          this.examples,
                                          (it) => it.name,
                                          (it, idx) => {
                                            return html`<option
                                              value="${idx}"
                                              title="${it.name}"
                                            >
                                              ${it.name}
                                            </option>`;
                                          }
                                        )
                                      }
                                    </select>
                                    </div>`}
          </div>
        </div>
        <div class="code-editor-container">
          <div class="wrapper" id="otlp-data-input"></div>
        </div>
      </div>
    `;
  }

  _handleExampleChanged(event) {
    if (!event.target.value) return;
    let idx = parseInt(event.target.value);
    let example = this.examples[idx];
    if (!example) return;

    this.payload = JSON.stringify(JSON.parse(example.value), null, 2);

    this.dispatchEvent(
      new CustomEvent('payload-example-changed', {
        detail: {value: example},
        bubbles: true,
        composed: true,
        cancelable: true,
      })
    );
  }

  _notifyPayloadChange(value) {
    this.dispatchEvent(
      new CustomEvent('payload-changed', {
        detail: {value: value},
        bubbles: true,
        composed: true,
      })
    );
  }

  _initCodeEditor() {
    this._editor = new EditorView({
      extensions: [
        basicSetup,
        Prec.highest(
          keymap.of([
            indentWithTab,
            {key: 'Enter', run: insertNewlineAndIndent, shift: () => true},
          ])
        ),
        linter(jsonParseLinter()),
        lintGutter(),
        this._editorReadOnlyCompartment.of(
          EditorState.readOnly.of(this.readOnly || false)
        ),
        EditorView.lineWrapping,
        EditorView.updateListener.of((v) => {
          if (v.docChanged) {
            this._notifyPayloadChange(this.payload);
          }
        }),
        json(),
      ],
      parent: this.shadowRoot.querySelector('#otlp-data-input'),
    });
  }
}

customElements.define('playground-payload-panel', PlaygroundPayloadPanel);
