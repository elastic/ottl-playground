import {html, LitElement} from 'lit-element';
import '../wasm_exec.js';
import Split from 'split.js';
import {basicSetup, EditorView} from 'codemirror';
import {json} from '@codemirror/lang-json';
import {PAYLOAD_EXAMPLES} from './examples';
import * as jsondiffpatch from 'jsondiffpatch';
import * as annotatedFormatter from 'jsondiffpatch/formatters/annotated';
import * as htmlFormatter from 'jsondiffpatch/formatters/html';
import './code-panels/config-panel';
import './code-panels/payload-panel';
import {playgroundStyles} from './playground.styles';

export class Playground extends LitElement {
  static properties = {
    // attributes
    title: {type: String, attribute: true},
    hideEvaluator: {type: Boolean, attribute: 'hide-evaluator'},
    hideRunButton: {type: Boolean, attribute: 'hide-run-button'},
    // states
    _loading: {type: Boolean, state: true},
    _evaluator: {type: String, state: true},
  };

  static get styles() {
    return playgroundStyles;
  }

  constructor() {
    super();
    this._initDefaultValues();
    this._addEventListeners();
  }

  firstUpdated() {
    this._loading = false;
    this._fetchWebAssembly();
  }

  _initDefaultValues() {
    this._loading = true;
    this.hideEvaluator = false;
    this.hideRunButton = false;
    this._evaluator = 'transform_processor';
    this.title = 'OTTL Playground';
  }

  _getPayloadPanelElement() {
    return this.shadowRoot.querySelector('#payload-code-panel');
  }

  _getConfigPanelElement() {
    return this.shadowRoot.querySelector('#config-code-panel');
  }

  get state() {
    return {
      config: this._getConfigPanelElement().text,
      payload: this._getPayloadPanelElement().text,
      payloadType: this._getPayloadPanelElement().payloadType,
      evaluator: this._evaluator,
    };
  }

  set state(state) {
    this._evaluator = state.evaluator;
    this._getConfigPanelElement().text = state.config;
    this._setPayloadValue(state.payloadType, state.payload);

    const event = new CustomEvent('playground-evaluator-change', {
      detail: {value: state.evaluator},
      bubbles: true,
      composed: true,
      cancelable: true,
    });

    this.dispatchEvent(event);
  }

  render() {
    return html`
            <div id="loading" style="display: ${this._loading ? '' : 'none'}">
                <svg xmlns="http://www.w3.org/2000/svg" width="36" height="36"
                     viewBox="0 0 24 24">
                    <style>.spinner_qM83 {
                        animation: spinner_8HQG 1.05s infinite
                    }

                    .spinner_oXPr {
                        animation-delay: .1s
                    }

                    .spinner_ZTLf {
                        animation-delay: .2s
                    }

                    @keyframes spinner_8HQG {
                        0%, 57.14% {
                            animation-timing-function: cubic-bezier(0.33, .66, .66, 1);
                            transform: translate(0)
                        }
                        28.57% {
                            animation-timing-function: cubic-bezier(0.33, 0, .66, .33);
                            transform: translateY(-6px)
                        }
                        100% {
                            transform: translate(0)
                        }
                    }</style>
                    <circle class="spinner_qM83" cx="4" cy="12" r="3"/>
                    <circle class="spinner_qM83 spinner_oXPr" cx="12" cy="12" r="3"/>
                    <circle class="spinner_qM83 spinner_ZTLf" cx="20" cy="12" r="3"/>
                </svg>
                </svg>
            </div>
            <div class="playground" id="playground"
                 style="display: ${this._loading ? 'none' : ''}">
                <slot name="playground-controls">
                    <playground-controls id="playground-controls"
                                         ?hide-run-button="${this.hideRunButton}"
                                         ?hide-evaluator=${this.hideEvaluator}
                                         evaluator="${this._evaluator}">
                        <slot name="playground-controls-app-title-text"
                              slot="app-title-text">
                            ${this.title}
                        </slot>
                        <slot name="playground-controls-custom-components"
                              slot="custom-components">
                        </slot>
                    </playground-controls>
                </slot>
                <div class="split-horizontal">
                    <div id="left-panel">
                        <div class="split-vertical">
                            <div id="config-code-panel-container">
                                <playground-config-panel id="config-code-panel"
                                                         evaluator="${this._evaluator}"
                                                         @playground-load-example="${this._handleLoadExample}">
                                    >
                                </playground-config-panel>
                            </div>
                            <div id="payload-code-panel-container">
                                <playground-payload-panel 
                                    id="payload-code-panel" 
                                    payload-type="logs"
                                >
                                </playground-payload-panel>
                            </div>
                        </div>
                    </div>
                    <div class="hidden-overflow" id="right-panel">
                        <div class="full-size">
                            <div class="result-panel-controls">
                                <div class="header">
                                    <span><strong>Result</strong> <sup><small>DIFF</small></sup></span>
                                </div>
                            </div>
                            <div>
                                <div class="result-panel-delta">
                                    <div>
                                        Delta
                                    </div>
                                    <div>
                                        <select class="delta-select"
                                                id="diff-delta-select"
                                                @change="${this._selectedDeltaChanged}">
                                            <option value="visual">Visual</option>
                                            <option value="json">JSON</option>
                                            <option value="annotated_json">Annotated
                                                JSON
                                            </option>
                                            <option value="logs">Execution Logs</option>
                                        </select>
                                    </div>
                                    <div id="show-unchanged-group"
                                         style="display: flex; align-items: center; cursor: default">
                                        <input id="show-unchanged-input" type="checkbox"
                                               @change="${this._showUnchangedInputChanged}">
                                        <div @click="${this._showUnchangedInputChanged}">
                                            Show unchanged
                                        </div>
                                        </input>
                                    </div>
                                </div>
                            </div>
                            <div class="result-panel-content" id="result-panel">
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
  }

  _addEventListeners() {
    window.addEventListener('load', () => {
      // Invoked after load, otherwise it won't respect the max heights
      this._spitComponents();
    });

    this.addEventListener('playground-run', () => {
      this._runStatements();
    });

    this.addEventListener('playground-evaluator-change', (e) => {
      this._evaluator = e.detail.value;
    });
  }

  _fetchWebAssembly() {
    // eslint-disable-next-line no-undef
    const go = new Go();
    WebAssembly.instantiateStreaming(
      fetch('ottlplayground.wasm'),
      go.importObject
    ).then((result) => {
      go.run(result.instance);
      this.updateComplete.then(() => {
        const event = new CustomEvent('playground-wasm-ready', {
          bubbles: true,
          composed: true,
          cancelable: true,
        });
        this.dispatchEvent(event);
      });
    });
  }

  _spitComponents() {
    Split(
      [
        this.shadowRoot.querySelector('#config-code-panel-container'),
        this.shadowRoot.querySelector('#payload-code-panel-container'),
      ],
      {
        direction: 'vertical',
      }
    );

    Split([
      this.shadowRoot.querySelector('#left-panel'),
      this.shadowRoot.querySelector('#right-panel'),
    ]);
  }

  _setEditorValue(editor, value) {
    editor.dispatch({
      changes: {from: 0, to: editor.state.doc.length, insert: value},
    });
  }

  _setResultText(text, pre = false) {
    let escaped = text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');

    let value = pre ? `<pre>${escaped}</pre>` : escaped;
    this.shadowRoot.querySelector('#result-panel').innerHTML =
      `<div class="text">${value}</div>`;
  }

  _setResultError(message) {
    delete window.lastRunData;
    this.shadowRoot.querySelector('#diff-delta-select').value = 'visual';
    this._setResultText(`Error executing statements: ${message}`);
  }

  _setPayloadExample(type) {
    let value = JSON.stringify(JSON.parse(PAYLOAD_EXAMPLES[type]), null, 2);
    this._setPayloadValue(type, value);
  }

  _setPayloadValue(type, value) {
    let el = this._getPayloadPanelElement();
    el.text = value;
    el.payloadType = type;
  }

  _runStatements() {
    let state = this.state;
    try {
      JSON.parse(state.payload);
    } catch (e) {
      this._setResultError(`Invalid OTLP JSON payload: ${e}`);
      return;
    }

    // eslint-disable-next-line no-undef
    let result = executeStatements(
      state.config,
      state.payloadType,
      state.payload,
      state.evaluator
    );
    if (result && Object.prototype.hasOwnProperty.call(result, 'error')) {
      this._setResultError(result.error);
      console.error('WASM error: ', result);
      return;
    }

    window.lastRunData = {
      input: state.payload,
      result: result,
    };

    this._setResult(state.payload, result);
  }

  _setLogResult(result) {
    let resultPanel = this.shadowRoot.querySelector('#result-panel');
    resultPanel.innerHTML = '';
    let editor = new EditorView({
      extensions: [basicSetup, EditorView.editable.of(false), json()],
      parent: resultPanel,
    });

    this._setEditorValue(editor, result.logs);
  }

  _setJsonResult(input, result) {
    if (!result.value) {
      this._setResultError('Empty result value');
      return;
    }

    let left = JSON.parse(input);
    let right = JSON.parse(result.value);
    let selectedDelta =
      this.shadowRoot.querySelector('#diff-delta-select').value;
    let resultPanel = this.shadowRoot.querySelector('#result-panel');

    // Raw JSON result
    if (selectedDelta === 'json') {
      resultPanel.innerHTML = '';
      let editor = new EditorView({
        extensions: [
          basicSetup,
          EditorView.editable.of(false),
          EditorView.lineWrapping,
          json(),
        ],
        parent: resultPanel,
      });
      this._setEditorValue(editor, JSON.stringify(right, null, 2));
      return;
    }

    // Comparable JSON results
    const delta = jsondiffpatch.diff(left, right);
    if (!delta) {
      this._setResultText('No changes');
      return;
    }

    let formatter =
      selectedDelta === 'annotated_json' ? annotatedFormatter : htmlFormatter;

    if (formatter.showUnchanged && formatter.hideUnchanged) {
      formatter.showUnchanged(
        this.shadowRoot.querySelector('#show-unchanged-input').checked,
        this.shadowRoot.querySelector('#result-panel')
      );
    }

    resultPanel.innerHTML = formatter.format(delta, left);
  }

  _setResult(input, result) {
    if (this.shadowRoot.querySelector('#diff-delta-select').value === 'logs') {
      this._setLogResult(result);
      return;
    }
    this._setJsonResult(input, result);
  }

  _showUnchangedInputChanged(e) {
    let el = this.shadowRoot.querySelector('#show-unchanged-input');
    if (el.disabled) {
      return;
    }
    if (typeof e.target.checked === 'undefined') {
      el.checked = !el.checked;
    }
    htmlFormatter.showUnchanged(
      el.checked,
      this.shadowRoot.querySelector('#result-panel')
    );
  }

  _selectedDeltaChanged(event) {
    let delta = event.target.value;
    let showUnchangedInput = this.shadowRoot.querySelector(
      '#show-unchanged-input'
    );
    showUnchangedInput.disabled = delta !== 'visual';
    showUnchangedInput.checked = false;

    let data = window.lastRunData;
    if (data) {
      this._setResult(data.input, data.result);
    }
  }

  _handleLoadExample(event) {
    let example = event.detail.example;
    if (example) {
      this._getConfigPanelElement().text = example.config;
      this._setPayloadExample(example.otlp_type);
    }
  }
}

customElements.define('playground-stage', Playground);
