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
import {globalStyles} from '../../styles.js';
import {repeat} from 'lit/directives/repeat.js';
import {nothing} from 'lit';
import './copy-link-button';

export class PlaygroundControls extends LitElement {
  static properties = {
    title: {type: String},
    executor: {type: String},
    executors: {type: Object},
    version: {type: String},
    versions: {type: Object},
    hideExecutors: {type: Boolean, attribute: 'hide-executors'},
    hideRunButton: {type: Boolean, attribute: 'hide-run-button'},
    hideCopyLinkButton: {type: Boolean, attribute: 'hide-copy-link-button'},
    loading: {type: Boolean},
  };

  constructor() {
    super();
    this.hideExecutors = false;
    this.hideRunButton = false;
    this.hideCopyLinkButton = false;
    this.loading = false;
    this.executor = 'transform_processor';
    this.versions = [{version: '-', artifact: 'ottlplayground.wasm'}];

    window.addEventListener('keydown', (event) => {
      if (event.shiftKey && event.key === 'Enter') {
        this._notifyRunRequested();
      }
    });
  }

  static get styles() {
    return [
      css`
        .playground-controls {
          overflow: hidden;
          width: 100%;
        }

        .playground-controls div {
          float: left;
          text-align: center;
          margin: 5px 0 5px 0;
          text-decoration: none;
        }

        .playground-controls .app-title {
          font-size: 25px;
          padding: 10px;
          font-weight: 600;
        }

        .playground-controls div.right {
          float: right;
          display: flex;
          align-items: center;
          justify-content: center;

          div:not(:last-child) {
            margin-right: 4px;
          }
        }

        .run-button {
          background-color: #04aa6d;
          border: 1px solid #049d65;
          color: white;
          padding: 10px;
          width: 75px;
          text-align: center;
          text-decoration: none;
          display: inline-block;
          font-size: 16px;
          cursor: pointer;
          border-radius: 4px;
        }

        .run-button:hover {
          background-color: #3e8e41;
        }

        #executor {
          width: 180px;
          height: 40px;
        }

        #version {
          height: 40px;
        }
      `,
      globalStyles,
    ];
  }

  render() {
    return html`
      <div class="playground-controls">
        <slot name="app-title">
          <div class="app-title">
            <slot name="app-title-text"><span>${this.title}</span></slot>
          </div>
        </slot>
        <div class="right">
          ${this.hideExecutors
            ? nothing
            : html`
                <div class="executor-container">
                  <select
                    id="version"
                    ?disabled="${this.loading}"
                    name="version"
                    @change="${this._notifyVersionChanged}"
                    title="Version"
                  >
                    <optgroup label="OpenTelemetry Collector Contrib">
                      ${this.versions &&
                      repeat(
                        this.versions,
                        (it) => it.version,
                        (it) => {
                          return html` <option
                            title="opentelemetry-collector-contrib (${it.version})"
                            ?selected="${it.version === this.version}"
                            value="${it.version}"
                          >
                            ${it.version}
                          </option>`;
                        }
                      )}
                    </optgroup>
                  </select>
                </div>
                <div class="executor-container">
                  <select
                    id="executor"
                    ?disabled="${this.loading}"
                    name="executor"
                    @change="${this._notifyExecutorChanged}"
                    title="Executor"
                  >
                    ${this.executors &&
                    repeat(
                      Object.entries(this._groupedExecutors()),
                      ([group]) => group,
                      ([group, executors]) => {
                        return html`<optgroup label="${group}">
                          ${repeat(
                            executors,
                            (it) => it.id,
                            (it) => {
                              return html` <option
                                title="${it.name} ${it.type} (${it.version})"
                                ?selected="${it.id === this.executor}"
                                value="${it.id}"
                              >
                                ${it.name}
                                ${it.id === this.executor ? it.type : ''}
                              </option>`;
                            }
                          )}
                        </optgroup>`;
                      }
                    )}
                  </select>
                </div>
              `}
          ${this.hideRunButton
            ? nothing
            : html`
                <div>
                  <button
                    class="run-button"
                    ?disabled="${this.loading}"
                    id="btn-run"
                    @click="${this._notifyRunRequested}"
                  >
                    <span id="run-btn">
                      ${this.loading
                        ? html`
                            <span>
                              <!-- prettier-ignore -->
                              <svg height="8px" id="Layer_1" style="enable-background: new 0 0 30 30;" viewBox="0 0 18 10" width="10px" x="0px" xml:space="preserve" xmlns="http://www.w3.org/2000/svg"><rect fill="#fff" height="20" width="4" x="0" y="0"><animate attributeName="opacity" attributeType="XML" begin="0s" dur="0.6s" repeatCount="indefinite" values="1; .2; 1"></animate></rect><rect fill="#fff" height="20" width="4" x="7" y="0"><animate attributeName="opacity" attributeType="XML" begin="0.2s" dur="0.6s" repeatCount="indefinite" values="1; .2; 1"></animate></rect><rect fill="#fff" height="20" width="4" x="14" y="0"><animate attributeName="opacity" attributeType="XML" begin="0.4s" dur="0.6s" repeatCount="indefinite" values="1; .2; 1"></animate></rect></svg>
                            </span>
                          `
                        : html`
                            <span class="tooltip"
                              >Run &#x25BA;
                              <span
                                class="tooltip-text tooltip-text-position-bottom"
                                >&#8679;+&#8629;</span
                              >
                            </span>
                          `}
                    </span>
                  </button>
                </div>
              `}
          ${this.hideCopyLinkButton
            ? nothing
            : html`
                <playground-copy-link-button></playground-copy-link-button>
              `}
          <slot name="custom-components"></slot>
        </div>
      </div>
    `;
  }

  _groupedExecutors() {
    const grouped = {};
    this.executors?.forEach((executor) => {
      let groupName =
        executor.type.charAt(0).toUpperCase() + executor.type.slice(1);
      if (!grouped[groupName]) {
        grouped[groupName] = [];
      }
      grouped[groupName].push(executor);
    });
    return grouped;
  }

  _notifyRunRequested() {
    this.dispatchEvent(
      new Event('playground-run-requested', {bubbles: true, composed: true})
    );
  }

  _notifyVersionChanged(e) {
    const event = new CustomEvent('version-changed', {
      detail: {
        value: e.target.value,
      },
      bubbles: true,
      composed: true,
    });
    this.dispatchEvent(event);
  }

  _notifyExecutorChanged(e) {
    const event = new CustomEvent('executor-changed', {
      detail: {
        value: e.target.value,
      },
      bubbles: true,
      composed: true,
    });
    this.dispatchEvent(event);
  }
}

customElements.define('playground-controls', PlaygroundControls);
