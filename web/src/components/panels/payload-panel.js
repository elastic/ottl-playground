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
import {keymap} from '@codemirror/view';
import {indentWithTab, insertNewlineAndIndent} from '@codemirror/commands';
import {PAYLOAD_EXAMPLES} from '../examples';
import {linter, lintGutter} from '@codemirror/lint';
import {json, jsonParseLinter} from '@codemirror/lang-json';
import {nothing} from 'lit';
import {Prec} from '@codemirror/state';

export class PlaygroundPayloadPanel extends LitElement {
  static properties = {
    payload: {type: String},
    hideExamples: {type: Boolean, attribute: 'hide-examples'},
    _editor: {type: Object, state: true},
  };

  constructor() {
    super();
    this.payload = '{}';
    this.hideLoadExample = false;
  }

  static get styles() {
    return [
      css`
        .example-button {
          background-color: #f1f1f1;
          border: none;
          color: black;
          font-size: 28px;
          padding: 0 12px;
          text-align: center;
          text-decoration: none;
          display: inline-block;
          margin: 4px 2px;
          cursor: pointer;
          border-radius: 5px;
        }
      `,
      ...codePanelsStyles,
    ];
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

  set selectedExample(val) {
    this.shadowRoot.querySelector('#example-select').value = val;
  }

  firstUpdated() {
    this._initCodeEditor();
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
            ${this.hideLoadExample
              ? nothing
              : html`
                                    <select
                                            id="example-select"
                                            @change="${this._handleExampleChanged}"
                                            title="Payload example">
                                        <option selected disabled value="">Example
                                        </option>
                                        <option value="logs">Logs</option>
                                        <option value="traces">Traces</option>
                                        <option value="metrics">Metrics</option>
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

  _handleExampleChanged(e) {
    if (!e.target.value) return;
    let val = JSON.stringify(
      JSON.parse(PAYLOAD_EXAMPLES[e.target.value]),
      null,
      2
    );
    this._editor?.dispatch({
      changes: {from: 0, to: this._editor.state.doc.length, insert: val},
    });
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
