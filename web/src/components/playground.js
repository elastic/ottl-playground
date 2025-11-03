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
import {DEFAULT_PAYLOAD_EXAMPLES, DEFAULT_PAYLOADS} from './examples';
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
    executor: {type: String},
    hideExecutors: {type: Boolean, attribute: 'hide-executors'},
    hideRunButton: {type: Boolean, attribute: 'hide-run-button'},
    disableShareLink: {type: Boolean, attribute: 'disable-share-link'},
    baseUrl: {type: String, attribute: 'base-url'},

    _loading: {state: true},
    _loadingWasm: {state: true},
    _executors: {state: true},
    _executor: {state: true},
    _versions: {state: true},
    _result: {state: true},
    _activeResult: {state: true},
    _debuggingInfo: {state: true},
  };

  constructor() {
    super();
    this._initDefaultValues();
    this._addEventListeners();
  }

  _initDefaultValues() {
    this._loading = true;
    this._hideExecutors = false;
    this.hideRunButton = false;
    this.disableShareLink = false;
    this.payload = '{}';
    this.baseUrl = '';
    this.executor = 'transform_processor';
  }

  static get styles() {
    return playgroundStyles;
  }

  get state() {
    return {
      version: this.version,
      executor: this.executor,
      payload: this.payload,
      config: this.config,
    };
  }

  set state(state) {
    this.version = state.version;
    this.executor = state.executor;
    this.payload = state.payload;
    this.config = state.config;
    // Reset the payload example dropdown
    this._clearSelectedPayloadExample();
  }

  firstUpdated() {
    this._splitComponents();
    this._loading = false;
    this._fetchWebAssemblyVersions().then(() => {
      this._initState();
    });
  }

  willUpdate(changedProperties) {
    if (
      changedProperties.has('_executors') ||
      changedProperties.has('executor')
    ) {
      this._executor = this._executors?.find((p) => p.id === this.executor);
    }
    super.willUpdate(changedProperties);
  }

  updated(changedProperties) {
    if (changedProperties.has('executor')) {
      this._debuggingInfo = null;
    }
    super.updated(changedProperties);
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
        // older versions uses "evaluator" instead of "executor"
        executor: urlStateData?.executor || urlStateData?.evaluator,
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
    let hash = this._getUrlHash();
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

  _clearSelectedPayloadExample() {
    let panel = this.shadowRoot.querySelector('#payload-code-panel');
    if (panel) {
      panel.clearSelectedExample();
    }
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
            ?hide-executors=${this._hideExecutors}
            ?hide-copy-link-button="${this.disableShareLink}"
            ?loading="${this._loadingWasm}"
            executor="${this.executor}"
            executors="${JSON.stringify(this._executors)}"
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
                  examples="${JSON.stringify(this._getConfigExamples())}"
                  config="${this.config}"
                  @config-changed="${(e) => (this.config = e.detail.value)}"
                  config-docs-url="${this._executor?.docsURL}"
                  @config-example-changed="${this._handleConfigExampleChanged}"
                  ?debugger-enabled="${this._executor?.debuggable === true}"
                  debugging-info="${JSON.stringify(this._debuggingInfo)}"
                  @debugging-line-changed="${this._handleDebuggingLineChanged}"
                  @debugging-stop-requested="${this
                    ._handDebuggingStopRequested}"
                >
                </playground-config-panel>
              </div>
              <div id="payload-code-panel-container">
                <playground-payload-panel
                  id="payload-code-panel"
                  payload="${this.payload}"
                  examples="${JSON.stringify(this._getPayloadExamples())}"
                  @payload-changed="${(e) => (this.payload = e.detail.value)}"
                  ?read-only="${this._debuggingInfo?.debugging === true}"
                >
                </playground-payload-panel>
              </div>
            </div>
          </div>
          <div class="hidden-overflow" id="right-panel">
            <playground-result-panel
              id="result-panel"
              payload="${this.payload}"
              result="${JSON.stringify(this._activeResult)}"
              view-config="${JSON.stringify(
                this._executors?.find((p) => p.id === this.executor)
                  ?.resultViewConfig || {}
              )}"
            >
            </playground-result-panel>
          </div>
        </div>
      </div>
    `;
  }

  _addEventListeners() {
    let me = this;
    window.wasmPanicHandler = (error) => {
      me._getResultPanel().showErrorMessage(error);
    };

    window.addEventListener('playground-wasm-ready', () => {
      // eslint-disable-next-line no-undef
      this._executors = getExecutors();
      if (!this._executors) {
        this.executor = '';
        this._executor = null;
      } else {
        if (!this.executor) {
          this.executor = this._executors[0]?.id;
          this._executor = this._executors[0];
        } else {
          let exec = this._executors.find((e) => e.id === this.executor);
          if (!exec) {
            this.executor = this._executors[0]?.id;
            this._executor = this._executors[0];
          }
        }
      }
    });

    this.addEventListener('playground-run-requested', () => {
      this._runStatements();
    });

    this.addEventListener('executor-changed', (e) => {
      this.executor = e.detail.value;
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

  _splitComponents() {
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
      this._getResultPanel().showErrorMessage(
        `Invalid OTLP JSON payload: ${e}`
      );
      return;
    }

    let debug = this.shadowRoot
      ?.querySelector('#config-code-panel')
      ?.hasBreakpoints();

    // eslint-disable-next-line no-undef
    let result = execute(
      state.config,
      payloadType,
      state.payload,
      state.executor,
      debug
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
    this._setResult(result);
  }

  _setResult(result) {
    this._result = result;
    if (result?.debug === true) {
      this._getResultPanel()
        .clearResult()
        .then(() => {
          this._setDebuggingInfo(result);
        });
    } else {
      if (this._debuggingInfo) {
        this._debuggingInfo = {...this._debuggingInfo, debugging: false};
      }
      this._activeResult = result;
    }
  }

  _setDebuggingInfo(result) {
    let debuggingInfo = {
      debugging: false,
      lines: [],
      lineResultIndex: {},
    };

    if (result?.debug === true) {
      let results;
      try {
        results = JSON.parse(result.value);
        result.value = results;
      } catch (e) {
        this._debuggingInfo = {...this._debuggingInfo, ...debuggingInfo};
        return;
      }

      debuggingInfo.debugging = true;
      for (let i = 0; i < results.length; i++) {
        let res = results[i];
        debuggingInfo.lines.push(res.line);
        debuggingInfo.lineResultIndex[res.line] = i;
      }
    }

    // If we were already in debugging mode, reset the debugging flag to false
    // so the UI can reset properly.
    if (this._debuggingInfo?.debugging === true) {
      this._debuggingInfo = null;
    }

    this.updateComplete.then(() => {
      this._debuggingInfo = {...this._debuggingInfo, ...debuggingInfo};
    });
  }

  _handleConfigExampleChanged(event) {
    let example = event.detail.value;
    if (example) {
      let payload = example.payload || DEFAULT_PAYLOADS[example.signal];
      this.payload = JSON.stringify(JSON.parse(payload), null, 2);
      this._clearSelectedPayloadExample();
    }
  }

  _handDebuggingStopRequested() {
    this._debuggingInfo = null;
  }

  _handleDebuggingLineChanged(event) {
    if (this._debuggingInfo) {
      let line = event.detail.value;
      if (line === -1) {
        this._getResultPanel()
          .clearResult()
          .then(() => {
            this._activeResult = null;
          });
      } else {
        let resultIndex = this._debuggingInfo.lineResultIndex[line];
        this._activeResult = this._result?.value?.[resultIndex] || null;
      }
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
    let urlHash = this._getUrlHash();
    if (urlHash) {
      return this._replaceUrlHash(value);
    } else {
      return this._getCurrentUrl() + '#' + value;
    }
  }

  _isEmbedded() {
    return window.self !== window.top;
  }

  _getUrlHash() {
    if (this._isEmbedded()) {
      return window.location.hash?.substring(1);
    } else {
      return window.top.location.hash?.substring(1);
    }
  }

  _getCurrentUrl() {
    if (this._isEmbedded()) {
      return window.location.href;
    } else {
      return window.top.location.href;
    }
  }

  _replaceUrlHash(value) {
    if (this._isEmbedded()) {
      return window.location.href.replace(window.location.hash, '#' + value);
    } else {
      return window.top.location.href.replace(
        window.top.location.hash,
        '#' + value
      );
    }
  }

  _sortBy(property) {
    return function (a, b) {
      return a[property] < b[property] ? -1 : a[property] > b[property] ? 1 : 0;
    };
  }

  _getConfigExamples() {
    if (!this._executor?.examples?.configs) {
      return null;
    }
    return this._executor?.examples?.configs.sort(this._sortBy('name'));
  }

  _getPayloadExamples() {
    if (!this._executor?.examples?.payloads) {
      return DEFAULT_PAYLOAD_EXAMPLES;
    }
    return this._executor?.examples?.payloads.sort(this._sortBy('name'));
  }

  _getResultPanel() {
    return this.shadowRoot.querySelector('#result-panel');
  }

  async _copyToClipboard(textToCopy) {
    if (navigator.clipboard && window.isSecureContext && !this._isEmbedded()) {
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
