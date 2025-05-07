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
import '../wasm_exec.js';
import Split from 'split.js';
import {CONFIG_EXAMPLES, PAYLOAD_EXAMPLES} from './examples';
import './panels/config-panel';
import './panels/payload-panel';
import './panels/result-panel';
import {playgroundStyles} from './playground.styles';
import {nothing} from 'lit';
import {getJsonPayloadType} from './utils/json-payload';
import {base64ToUtf8, utf8ToBase64} from './utils/base64';

export class Playground extends LitElement {
  static properties = {
    title: {type: String},
    config: {type: String},
    payload: {type: String},
    version: {type: String},
    evaluator: {type: String},
    hideEvaluators: {type: Boolean, attribute: 'hide-evaluators'},
    hideRunButton: {type: Boolean, attribute: 'hide-run-button'},
    disableShareLink: {type: Boolean, attribute: 'disable-share-link'},
    baseUrl: {type: String, attribute: 'base-url'},

    _loading: {state: true},
    _loadingWasm: {state: true},
    _evaluators: {state: true},
    _evaluatorsDocsURL: {state: true},
    _versions: {state: true},
    _result: {state: true},
  };

  constructor() {
    super();
    this._initDefaultValues();
    this._addEventListeners();
  }

  _initDefaultValues() {
    this._loading = true;
    this._hideEvaluators = false;
    this.hideRunButton = false;
    this.disableShareLink = false;
    this.payload = '{}';
    this.baseUrl = '';
  }

  static get styles() {
    return playgroundStyles;
  }

  get state() {
    return {
      version: this.version,
      evaluator: this.evaluator,
      payload: this.payload,
      config: this.config,
    };
  }

  set state(state) {
    this.version = state.version;
    this.evaluator = state.evaluator;
    this.payload = state.payload;
    this.config = state.config;
    // Reset the payload example dropdown
    this._setSelectedPayloadExample('');
  }

  firstUpdated() {
    this._spitComponents();
    this._loading = false;
    this._fetchWebAssemblyVersions().then(() => {
      this._initState();
    });
  }

  willUpdate(changedProperties) {
    if (changedProperties.has('_evaluators')) {
      this._computeEvaluatorsDocsURL();
    }
    super.willUpdate(changedProperties);
  }

  _initState() {
    let urlStateData = this._loadURLBase64DataHash();
    let version = urlStateData?.version;
    if (!version || !this._versions?.some((it) => it.version === version)) {
      version = this._versions?.[0]?.version;
    }

    if (urlStateData) {
      this.state = {
        version: version,
        evaluator: urlStateData?.evaluator,
        payload: urlStateData?.payload ?? '{}',
        config: urlStateData?.config,
      };
    }

    this._fetchWebAssembly(this._resolveWebAssemblyArtifact(version));
  }

  _fetchWebAssemblyVersions() {
    return fetch('wasm/versions.json')
      .then((response) => {
        return response.json();
      })
      .then((json) => {
        this.version = this.version || json.versions?.[0]?.version;
        this._versions = json.versions;
      });
  }

  _loadURLBase64DataHash() {
    if (this.disableShareLink === true) return;
    let hash = window.top.location.hash?.substring(1);
    if (hash) {
      try {
        let data = JSON.parse(base64ToUtf8(hash));
        if (data.payload) {
          try {
            data.payload = JSON.stringify(JSON.parse(data.payload), null, 2);
          } catch (e) {
            // Ignore
          }
        }
        return data;
      } catch (e) {
        return null;
      }
    }
  }

  _setSelectedPayloadExample(example) {
    let panel = this.shadowRoot.querySelector('#payload-code-panel');
    if (panel) {
      panel.selectedExample = example;
    }
  }

  _computeEvaluatorsDocsURL() {
    let docsURLs = {};
    this._evaluators?.forEach((it) => {
      docsURLs[it.id] = it.docsURL ?? null;
    });
    this._evaluatorsDocsURL = docsURLs;
  }

  render() {
    return html`
      ${this._loading
        ? html`
            <div id="loading">
              <!-- prettier-ignore -->
              <svg xmlns="http://www.w3.org/2000/svg" width="36" height="36" viewBox="0 0 24 24"> <style>.spinner_qM83 { animation: spinner_8HQG 1.05s infinite } .spinner_oXPr { animation-delay: .1s } .spinner_ZTLf { animation-delay: .2s } @keyframes spinner_8HQG { 0%, 57.14% { animation-timing-function: cubic-bezier(0.33, .66, .66, 1); transform: translate(0) } 28.57% { animation-timing-function: cubic-bezier(0.33, 0, .66, .33); transform: translateY(-6px) } 100% { transform: translate(0) } }</style> <circle class="spinner_qM83" cx="4" cy="12" r="3"/> <circle class="spinner_qM83 spinner_oXPr" cx="12" cy="12" r="3"/> <circle class="spinner_qM83 spinner_ZTLf" cx="20" cy="12" r="3"/> </svg>
            </div>
          `
        : nothing}
      <div class="playground" id="playground">
        <slot name="playground-controls">
          <playground-controls
            id="playground-controls"
            ?hide-run-button="${this.hideRunButton}"
            ?hide-evaluators=${this._hideEvaluators}
            ?hide-copy-link-button="${this.disableShareLink}"
            ?loading="${this._loadingWasm}"
            evaluator="${this.evaluator}"
            evaluators="${JSON.stringify(this._evaluators)}"
            version="${this.version}"
            versions="${JSON.stringify(this._versions)}"
            @copy-link-click="${this._handleCopyLinkClick}"
          >
            <slot
              name="playground-controls-app-title-text"
              slot="app-title-text"
            >
              ${this.title
                ? html` ${this.title}&nbsp;<sup
                      class="beta-box"
                      title="The OTTL Playground is still in beta and the authors of this tool would welcome your feedback"
                      >BETA</sup
                    >`
                : nothing}
            </slot>
            <slot
              name="playground-controls-custom-components"
              slot="custom-components"
            >
            </slot>
          </playground-controls>
        </slot>
        <div class="split-horizontal">
          <div id="left-panel">
            <div class="split-vertical">
              <div id="config-code-panel-container">
                <playground-config-panel
                  id="config-code-panel"
                  examples="${JSON.stringify(CONFIG_EXAMPLES[this.evaluator])}"
                  config="${this.config}"
                  @config-changed="${(e) => (this.config = e.detail.value)}"
                  config-docs-url="${this._evaluatorsDocsURL?.[this.evaluator]}"
                  @config-example-changed="${this._handleConfigExampleChanged}"
                >
                  >
                </playground-config-panel>
              </div>
              <div id="payload-code-panel-container">
                <playground-payload-panel
                  id="payload-code-panel"
                  payload="${this.payload}"
                  @payload-changed="${(e) => (this.payload = e.detail.value)}"
                >
                </playground-payload-panel>
              </div>
            </div>
          </div>
          <div class="hidden-overflow" id="right-panel">
            <playground-result-panel
              id="result-panel"
              payload="${this.payload}"
              result="${JSON.stringify(this._result)}"
            >
            </playground-result-panel>
          </div>
        </div>
      </div>
    `;
  }

  _addEventListeners() {
    window.addEventListener('playground-wasm-ready', () => {
      // eslint-disable-next-line no-undef
      this._evaluators = statementsExecutors();
      if (!this._evaluators) {
        this.evaluator = '';
      } else {
        if (!this.evaluator) {
          this.evaluator = this._evaluators[0]?.id;
        } else if (!this._evaluators.some((e) => e.id === this.evaluator)) {
          this.evaluator = this._evaluators[0]?.id;
        }
      }
    });

    this.addEventListener('playground-run-requested', () => {
      this._runStatements();
    });

    this.addEventListener('evaluator-changed', (e) => {
      this.evaluator = e.detail.value;
    });

    this.addEventListener('version-changed', (e) => {
      this.version = e.detail.value;
      this._fetchWebAssembly(this._resolveWebAssemblyArtifact(this.version));
    });
  }

  _resolveWebAssemblyArtifact(version) {
    return this._versions.find((e) => e.version === version)?.artifact;
  }

  _fetchWebAssembly(artifact) {
    // eslint-disable-next-line no-undef
    const go = new Go();
    this._loadingWasm = true;

    let wasmUrl = this.baseUrl
      ? new URL(artifact, this.baseUrl).href
      : artifact;

    WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then(
      (result) => {
        go.run(result.instance);
        this.updateComplete.then(() => {
          this._loadingWasm = false;
          const event = new CustomEvent('playground-wasm-ready', {
            detail: {
              value: artifact,
            },
            bubbles: true,
            composed: true,
            cancelable: true,
          });
          window.dispatchEvent(event);
        });
      }
    );
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

  _runStatements() {
    let state = this.state;
    let payloadType;
    try {
      payloadType = getJsonPayloadType(this.payload);
    } catch (e) {
      this.shadowRoot
        .querySelector('#result-panel')
        .showErrorMessage(`Invalid OTLP JSON payload: ${e}`);
      return;
    }

    // eslint-disable-next-line no-undef
    let result = executeStatements(
      state.config,
      payloadType,
      state.payload,
      state.evaluator
    );

    this.dispatchEvent(
      new CustomEvent('playground-run-result', {
        detail: {
          state: state,
          result: result,
          error:
            result && Object.prototype.hasOwnProperty.call(result, 'error'),
        },
        bubbles: true,
        composed: true,
        cancelable: true,
      })
    );

    this.payload = state.payload;
    this._result = result;
  }

  _handleConfigExampleChanged(event) {
    let example = event.detail.value;
    if (example) {
      let payload = example.payload || PAYLOAD_EXAMPLES[example.otlp_type];
      this.payload = JSON.stringify(JSON.parse(payload), null, 2);
      this._setSelectedPayloadExample(example.otlp_type);
    }
  }

  _handleCopyLinkClick() {
    let data = {...this.state};
    try {
      // Try to linearize the JSON to make it smaller
      data.payload = JSON.stringify(JSON.parse(data.payload));
    } catch (e) {
      // Ignore and use it as it's
    }

    let key = utf8ToBase64(JSON.stringify(data));
    this._copyToClipboard(this._buildUrlWithLink(key)).catch((e) => {
      console.error(e);
    });

    document.location.hash = key;
  }

  _buildUrlWithLink(value) {
    if (window.top.location.hash) {
      return window.top.location.href.replace(
        window.top.location.hash,
        '#' + value
      );
    } else {
      return window.top.location.href + '#' + value;
    }
  }

  async _copyToClipboard(textToCopy) {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(textToCopy);
    } else {
      const textArea = document.createElement('textarea');
      textArea.value = textToCopy;
      textArea.style.position = 'absolute';
      textArea.style.left = '-999999px';
      document.body.prepend(textArea);
      textArea.select();
      try {
        document.execCommand('copy');
      } catch (error) {
        console.error(error);
      } finally {
        textArea.remove();
      }
    }
  }
}

customElements.define('playground-stage', Playground);
