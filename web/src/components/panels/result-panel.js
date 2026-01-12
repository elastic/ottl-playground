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
import * as htmlFormatter from 'jsondiffpatch/formatters/html';
import {basicSetup, EditorView} from 'codemirror';
import {json} from '@codemirror/lang-json';
import * as jsondiffpatch from 'jsondiffpatch';
import * as annotatedFormatter from 'jsondiffpatch/formatters/annotated';
import {resultPanelStyles} from './result-panel.styles.js';
import {escapeHTML} from '../utils/escape-html';
import {repeat} from 'lit/directives/repeat.js';
import {nothing} from 'lit';
import {Compartment} from '@codemirror/state';

const VIEW_VISUAL_DELTA = 'visual_delta';
const VIEW_ANNOTATED_DELTA = 'annotated_delta';
const VIEW_JSON = 'json';
const VIEW_LOGS = 'logs';

export class PlaygroundResultPanel extends LitElement {
  static properties = {
    view: {type: String},
    viewConfig: {type: Object, attribute: 'view-config'},
    payload: {type: String},
    result: {type: Object},

    _wrapLines: {state: true},
    _showUnchanged: {state: true},
    _errored: {state: true},
    _views: {state: true},
  };

  constructor() {
    super();
    this.view = VIEW_VISUAL_DELTA;
    this.viewConfig = {
      [VIEW_VISUAL_DELTA]: {enabled: true},
      [VIEW_ANNOTATED_DELTA]: {enabled: true},
      [VIEW_JSON]: {enabled: true},
      [VIEW_LOGS]: {enabled: true},
    };
    this._wrapLines = false;
    this._showUnchanged = false;
    this._jsonViewEditor = null;
    this._logsViewEditor = null;
    this._wrapLinesCompartment = new Compartment();
    this._updateResultViewSelect();
  }

  static get styles() {
    return resultPanelStyles;
  }

  willUpdate(changedProperties) {
    if (
      changedProperties.has('viewConfig') &&
      changedProperties.get('viewConfig')
    ) {
      this._updateResultViewSelect();
    }
    super.willUpdate(changedProperties);
  }

  updated(changedProperties) {
    if (
      changedProperties.has('result') ||
      (changedProperties.has('view') && changedProperties.get('view'))
    ) {
      let rerender = this._isReRender(changedProperties);
      this._renderResult(rerender);
    }
    super.updated(changedProperties);
  }

  showErrorMessage(message) {
    this.result = null;
    this._errored = true;
    this._renderResultText(message);
  }

  clearResult() {
    this.result = null;
    this._errored = false;
    this._resultPanel().innerHTML = '';
    return this.updateComplete;
  }

  render() {
    return html`
      ${this._errored
        ? html`
            <style>
              .result-panel-controls {
                border-left: red 4px solid !important;
              }
            </style>
          `
        : nothing}
      <div class="full-size">
        <div class="result-panel-controls">
          <div class="header">
            <span><strong>Result</strong></span>
          </div>
          ${this.result?.executionTime !== undefined
            ? html`
                <div
                  class="execution-time-header"
                  title="Estimated execution time"
                >
                  <span>
                    <span>${this.result?.executionTime} ms</span>
                  </span>
                </div>
              `
            : nothing}
        </div>
        <div>
          <div class="result-panel-view">
            <div>View</div>
            <div>
              <select
                class="view-select"
                id="diff-view-select"
                .value="${this.view}"
                @change="${this._selectedViewChanged}"
              >
                ${this._views &&
                repeat(
                  this._views,
                  (it) => it.id,
                  (it) => {
                    return html` <option
                      title="${it.name}"
                      ?selected="${it.id === this.view}"
                      value="${it.id}"
                    >
                      ${it.name}
                    </option>`;
                  }
                )}
              </select>
            </div>
            ${this.view === VIEW_VISUAL_DELTA
              ? html`
                        <div id="show-unchanged-group" class="result-panel-flex-group">
                            <input id="show-unchanged-input" type="checkbox" ?checked="${this._showUnchanged}"
                                   @change="${this._showUnchangedInputChanged}">
                            <div @click="${this._showUnchangedInputChanged}">
                                Show unchanged
                            </div>
                            </input>
                        </div>`
              : nothing}
            ${this._showWrapLinesOption()
              ? html`
                        <div id="wrap-lines-group" class="result-panel-flex-group">
                          <input id="wrap-lines-input" type="checkbox" ?checked="${this._wrapLines}"
                                 @change="${this._wrapLinesInputChanged}">
                              <div @click="${this._wrapLinesInputChanged}">
                                Wrap lines
                              </div>
                          </input>
                        </div>`
              : nothing}
          </div>
        </div>
        <div class="result-panel-content" id="result-panel"></div>
      </div>
    `;
  }

  _isReRender(changedProperties) {
    // If no result rendered yet, need to render
    if (this._resultPanel().childElementCount === 0) {
      return false;
    }
    // If view didn't change, and it's currently not null, no re-render needed
    if (!changedProperties.has('view') && this.view !== null) {
      return true;
    }
    // same value, no re-render needed
    return !!(
      changedProperties.has('view') &&
      changedProperties.get('view') === this.view
    );
  }

  _showWrapLinesOption() {
    return this.view && (this.view === VIEW_JSON || this.view === VIEW_LOGS);
  }

  _selectedViewChanged(event) {
    this.view = event.target.value;
  }

  _showUnchangedInputChanged() {
    let el = this._showUnchangedInput();
    if (el.disabled) {
      return;
    }

    this._showUnchanged = !this._showUnchanged;
    el.checked = this._showUnchanged;
    htmlFormatter.showUnchanged(this._showUnchanged, this._resultPanel());
  }

  _wrapLinesInputChanged() {
    let el = this._wrapLinesInput();
    if (el.disabled) {
      return;
    }

    this._wrapLines = !this._wrapLines;
    el.checked = this._wrapLines;
    [this._jsonViewEditor, this._logsViewEditor].forEach((view) => {
      if (view) {
        view.dispatch({
          effects: this._wrapLinesCompartment.reconfigure(
            this._wrapLines ? EditorView.lineWrapping : []
          ),
        });
      }
    });
  }

  _resultPanel() {
    return this.shadowRoot.querySelector('#result-panel');
  }

  _showUnchangedInput() {
    return this.shadowRoot.querySelector('#show-unchanged-input');
  }

  _wrapLinesInput() {
    return this.shadowRoot?.querySelector('#wrap-lines-input');
  }

  _updateResultViewSelect() {
    if (this.viewConfig) {
      this._views = [];
      if (this.viewConfig[VIEW_VISUAL_DELTA]?.enabled) {
        this._views.push({id: VIEW_VISUAL_DELTA, name: 'Visual diff'});
      }
      if (this.viewConfig[VIEW_ANNOTATED_DELTA]?.enabled) {
        this._views.push({id: VIEW_ANNOTATED_DELTA, name: 'Annotated diff'});
      }
      if (this.viewConfig[VIEW_JSON]?.enabled) {
        this._views.push({id: VIEW_JSON, name: 'JSON'});
      }
      if (this.viewConfig[VIEW_LOGS]?.enabled) {
        this._views.push({id: VIEW_LOGS, name: 'Execution logs'});
      }
      if (this.viewConfig[this.view]?.enabled === false) {
        this.view = this._views[0].id;
      }
    }
  }

  _renderResult(rerender = false) {
    if (!this.result) return;

    if (rerender === false) {
      this._resultPanel().innerHTML = '';
    }

    this._errored = !!this.result.error;
    if (this.view === VIEW_LOGS) {
      this._renderExecutionLogsResult(rerender);
      return;
    }

    let resultError = this.result?.error;
    if (resultError) {
      this._renderResultText(resultError);
    } else {
      this._renderJsonDiffResult(this.result, rerender);
    }
  }

  _renderExecutionLogsResult(rerender) {
    if (!this._logsViewEditor) {
      const selectionHighlight = EditorView.theme({
        '&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection':
          {backgroundColor: 'rgb(255 255 0 / 50%)'},
      });

      let extensions = [
        basicSetup,
        selectionHighlight,
        EditorView.editable.of(false),
        json(),
      ];
      if (this._wrapLinesInput()?.checked) {
        extensions.push(this._wrapLinesCompartment.of(EditorView.lineWrapping));
      } else {
        extensions.push(this._wrapLinesCompartment.of([]));
      }

      this._resultPanel().innerHTML = '';
      this._logsViewEditor = new EditorView({
        extensions: extensions,
        parent: this._resultPanel(),
      });
    } else {
      if (rerender !== true) {
        this._resultPanel().appendChild(this._logsViewEditor.dom);
      }
    }

    let anchor = this._logsViewEditor.state?.selection?.main?.anchor;
    if (anchor > this.result.logs.length) {
      anchor = this.result.logs.length;
    }

    this._logsViewEditor.dispatch({
      changes: {
        from: 0,
        to: this._logsViewEditor.state.doc.length,
        insert: this.result.logs,
      },
      selection: {anchor: anchor},
    });
  }

  _renderJsonDiffResult(result, rerender) {
    if (!result.value) {
      this._renderResultText('Empty result');
      return;
    }

    let left = JSON.parse(this.payload);
    let right = JSON.parse(result.value);
    if (this.view === VIEW_JSON) {
      if (!this._jsonViewEditor) {
        let extensions = [basicSetup, EditorView.editable.of(false), json()];

        if (this._wrapLinesInput()?.checked) {
          extensions.push(
            this._wrapLinesCompartment.of(EditorView.lineWrapping)
          );
        } else {
          extensions.push(this._wrapLinesCompartment.of([]));
        }

        this._jsonViewEditor = new EditorView({
          extensions: extensions,
          parent: this._resultPanel(),
        });
      } else {
        if (rerender !== true) {
          this._resultPanel().appendChild(this._jsonViewEditor.dom);
        }
      }

      let value = JSON.stringify(right, null, 2);
      if (result.json) {
        value = JSON.stringify(JSON.parse(result.json), null, 2);
      }

      this._jsonViewEditor.dispatch({
        changes: {
          from: 0,
          to: this._jsonViewEditor.state.doc.length,
          insert: value,
        },
      });
      return;
    }

    // Comparable JSON results
    const delta = jsondiffpatch
      .create({
        objectHash: function (obj, index) {
          // Spans
          if (obj?.spanId && obj?.traceId) {
            return `${obj.spanId}-${obj?.traceId}`;
          }
          // Metrics
          if (
            obj?.name &&
            (obj?.gauge ||
              obj?.sum ||
              obj?.histogram ||
              obj?.exponentialHistogram ||
              obj?.summary)
          ) {
            return obj?.name;
          }
          // Attributes
          if (obj?.key && obj?.value) {
            return obj.key;
          }
          return '$$index:' + index;
        },
      })
      .diff(left, right);
    if (!delta) {
      this._renderResultText('No changes');
      return;
    }

    let formatter =
      this.view === VIEW_ANNOTATED_DELTA ? annotatedFormatter : htmlFormatter;

    if (formatter.showUnchanged && formatter.hideUnchanged) {
      formatter.showUnchanged(this._showUnchanged, this._resultPanel());
    }
    this._resultPanel().innerHTML = formatter.format(delta, left);
  }

  _renderResultText(text) {
    let resultPanel = this._resultPanel();
    resultPanel.innerHTML = `<div class="text">${escapeHTML(text)}</div>`;
  }
}

customElements.define('playground-result-panel', PlaygroundResultPanel);
