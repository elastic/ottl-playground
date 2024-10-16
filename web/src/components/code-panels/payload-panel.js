import {css, html, LitElement} from 'lit-element';
import {codePanelsStyles} from './styles';
import {basicSetup, EditorView} from 'codemirror';
import {keymap} from '@codemirror/view';
import {indentWithTab} from '@codemirror/commands';
import {PAYLOAD_EXAMPLES} from '../examples';
import {linter, lintGutter} from '@codemirror/lint';
import {json, jsonParseLinter} from '@codemirror/lang-json';
import {nothing} from "lit";

export class PlaygroundPayloadPanel extends LitElement {
  static properties = {
    // attributes
    payloadType: {type: String, attribute: 'payload-type', reflect: true},
    payload: {type: String, attribute: true},
    hideLoadExample: {type: Boolean, attribute: 'hide-load-example'},
    // states
    _editor: {type: Object, state: true},
  };

  constructor() {
    super();
    this.payloadType = 'logs';
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

  get text() {
    if (!this._editor) return '';
    return this._editor.state.doc.toString();
  }

  set text(value) {
    this._editor.dispatch({
      changes: {from: 0, to: this._editor.state.doc.length, insert: value},
    });
  }

  firstUpdated(_changedProperties) {
    this._initCodeEditor();
  }

  render() {
    return html`
      <div class="code-panel-parent" id="input-sample-panel">
        <div class="code-panel-controls">
          <div class="header">
            <span
              ><strong>OTLP payload</strong>
              <sup><small>JSON</small></sup></span
            >
          </div>
          <div class="right" style="display: flex">
            <select
              id="data-type-input"
              .value="${this.payloadType}"
              @change="${this._handleTypeChange}"
            >
              <option value="logs">Logs</option>
              <option value="traces">Traces</option>
              <option value="metrics">Metrics</option>
            </select>
            ${ this.hideLoadExample ? nothing : html `
            <button
              class="example-button"
              id="load-example-button"
              title="Load selected data-type example"
              @click="${this._loadPayloadExample}"
            >
              &#9112;
            </button>
            `}
          </div>
        </div>
        <div class="code-editor-container">
          <div class="wrapper" id="otlp-data-input"></div>
        </div>
      </div>
    `;
  }

  _handleTypeChange(e) {
    this.payloadType = e.target.value;
    const event = new CustomEvent('playground-payload-type-change', {
      bubbles: true,
      detail: {value: this.payloadType},
    });
    this.dispatchEvent(event);
  }

  _loadPayloadExample() {
    this.text = JSON.stringify(
      JSON.parse(PAYLOAD_EXAMPLES[this.payloadType]),
      null,
      2
    );
  }

  _initCodeEditor() {
    this._editor = new EditorView({
      extensions: [
        basicSetup,
        keymap.of([indentWithTab]),
        linter(jsonParseLinter()),
        lintGutter(),
        EditorView.lineWrapping,
        json(),
      ],
      parent: this.shadowRoot.querySelector('#otlp-data-input'),
    });

    if (this.payload) {
        this._editor.dispatch({
            changes: {from: 0, to:  this._editor.state.doc.length, insert: this.payload},
        });
    }
  }
}

customElements.define('playground-payload-panel', PlaygroundPayloadPanel);
