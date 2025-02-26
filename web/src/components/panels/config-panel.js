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
import {Prec} from '@codemirror/state';
import {keymap} from '@codemirror/view';
import {indentWithTab, insertNewlineAndIndent} from '@codemirror/commands';
import {yaml} from '@codemirror/lang-yaml';
import {nothing} from 'lit';
import {repeat} from 'lit/directives/repeat.js';

export class PlaygroundConfigPanel extends LitElement {
  static properties = {
    examples: {type: Array},
    hideExamples: {type: Boolean, attribute: 'hide-examples'},
    config: {type: String},
    configDocsURL: {type: String, attribute: 'config-docs-url'},
    _editor: {state: true},
  };

  constructor() {
    super();
    this.hideExamples = false;
    this.examples = [];
    this.configDocsURL = '';
  }

  static get styles() {
    return codePanelsStyles;
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
                    style="max-width: 250px"
                  >
                    <option selected disabled value="">Example</option>
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
        <div class="code-editor-container">
          <div class="wrapper" id="config-input"></div>
        </div>
      </div>
    `;
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

  _notifyConfigChange(value) {
    this.dispatchEvent(
      new CustomEvent('config-changed', {
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
        EditorView.lineWrapping,
        yaml(),
        EditorView.updateListener.of((v) => {
          if (v.docChanged) {
            this._notifyConfigChange(this.config);
          }
        }),
      ],
      parent: this.shadowRoot.querySelector('#config-input'),
    });
  }
}

customElements.define('playground-config-panel', PlaygroundConfigPanel);
