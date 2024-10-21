import {html, LitElement} from 'lit-element';
import {codePanelsStyles} from './styles';
import {basicSetup, EditorView} from 'codemirror';
import {keymap} from '@codemirror/view';
import {indentWithTab} from '@codemirror/commands';
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
                  >
                    <option selected disabled value="">Example</option>
                    ${this.examples &&
                    repeat(
                      this.examples,
                      (it) => it.name,
                      (it, idx) => {
                        return html`<option value="${idx}">${it.name}</option>`;
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
        keymap.of([indentWithTab]),
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
