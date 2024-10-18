import {html, LitElement} from 'lit-element';
import * as htmlFormatter from 'jsondiffpatch/formatters/html';
import {basicSetup, EditorView} from 'codemirror';
import {json} from '@codemirror/lang-json';
import * as jsondiffpatch from 'jsondiffpatch';
import * as annotatedFormatter from 'jsondiffpatch/formatters/annotated';
import {resultPanelStyles} from './result-panel.styles.js';
import {escapeHTML} from '../utils/escape-html';
import {nothing} from 'lit';

export class PlaygroundResultPanel extends LitElement {
  static properties = {
    delta: {type: String},
    payload: {type: String},
    result: {type: Object},
    _errored: {state: true},
  };

  constructor() {
    super();
    this.delta = 'visual';
  }

  static get styles() {
    return resultPanelStyles;
  }

  updated(changedProperties) {
    if (
      changedProperties.has('result') ||
      (changedProperties.has('delta') && changedProperties.get('delta'))
    ) {
      this._renderResult();
    }
    super.updated(changedProperties);
  }

  showResult(payload, result) {
    this.payload = payload;
    this.result = result;
    this._renderResult(payload, result);
  }

  showErrorMessage(message) {
    this.result = null;
    this._errored = true;
    this._renderResultText(message);
  }

  render() {
    return html`
            ${
              this._errored
                ? html`
                    <style>
                      .result-panel-controls {
                        border-left: red 4px solid !important;
                      }
                    </style>
                  `
                : nothing
            }
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
                                    .value="${this.delta}"
                                    @change="${this._selectedDeltaChanged}">
                                <option value="visual">Visual</option>
                                <option value="json">JSON</option>
                                <option value="annotated_json">Annotated
                                    JSON
                                </option>
                                <option value="logs">Execution Logs</option>
                            </select>
                        </div>
                        <div id="show-unchanged-group" class="result-panel-flex-group">
                            <input id="show-unchanged-input" type="checkbox"
                                   @change="${this._showUnchangedInputChanged}">
                            <div @click="${this._showUnchangedInputChanged}">
                                Show unchanged
                            </div>
                            </input>
                        </div>
                        <div id="wrap-lines-group" class="result-panel-flex-group" style="display: none">
                          <input id="wrap-lines-input" type="checkbox"
                                 @change="${this._wrapLinesInputChanged}">
                              <div @click="${this._wrapLinesInputChanged}">
                                Wrap lines
                              </div>
                          </input>
                        </div>
                    </div>
                </div>
                <div class="result-panel-content" id="result-panel">
                </div>
            </div>
        `;
  }

  _showWrapLinesOption() {
    return this.delta && (this.delta === 'json' || this.delta === 'logs');
  }

  _selectedDeltaChanged(event) {
    this.delta = event.target.value;

    this.shadowRoot.querySelector('#show-unchanged-group').style.display =
      this.delta !== 'visual' ? 'none' : '';

    this.shadowRoot.querySelector('#wrap-lines-group').style.display =
      this._showWrapLinesOption() ? '' : 'none';
  }

  _showUnchangedInputChanged(e) {
    let el = this._showUnchangedInput();
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

  _wrapLinesInputChanged(e) {
    let el = this._wrapLinesInput();
    if (el.disabled) {
      return;
    }
    if (typeof e.target.checked === 'undefined') {
      el.checked = !el.checked;
    }
    this._renderResult();
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

  _renderResult() {
    if (!this.result) return;

    this._resultPanel().innerHTML = '';
    this._errored = !!this.result.error;

    if (this.delta === 'logs') {
      this._renderExecutionLogsResult();
      return;
    }

    let resultError = this.result?.error;
    if (resultError) {
      this._renderResultText(resultError);
    } else {
      this._renderJsonDiffResult();
    }
  }

  _renderExecutionLogsResult() {
    let extensions = [basicSetup, EditorView.editable.of(false), json()];

    if (this._wrapLinesInput()?.checked) {
      extensions.push(EditorView.lineWrapping);
    }

    let editor = new EditorView({
      extensions: extensions,
      parent: this._resultPanel(),
    });

    editor.dispatch({
      changes: {
        from: 0,
        to: editor.state.doc.length,
        insert: this.result.logs,
      },
    });
  }

  _renderJsonDiffResult() {
    if (!this.result.value) {
      this._renderResultText('Empty result');
      return;
    }

    let left = JSON.parse(this.payload);
    let right = JSON.parse(this.result.value);
    if (this.delta === 'json') {
      let extensions = [basicSetup, EditorView.editable.of(false), json()];

      if (this._wrapLinesInput()?.checked) {
        extensions.push(EditorView.lineWrapping);
      }

      let editor = new EditorView({
        extensions: extensions,
        parent: this._resultPanel(),
      });

      editor.dispatch({
        changes: {
          from: 0,
          to: editor.state.doc.length,
          insert: JSON.stringify(right, null, 2),
        },
      });
      return;
    }

    // Comparable JSON results
    const delta = jsondiffpatch.diff(left, right);
    if (!delta) {
      this._renderResultText('No changes');
      return;
    }

    let formatter =
      this.delta === 'annotated_json' ? annotatedFormatter : htmlFormatter;

    if (formatter.showUnchanged && formatter.hideUnchanged) {
      formatter.showUnchanged(
        this._showUnchangedInput().checked,
        this._resultPanel()
      );
    }
    this._resultPanel().innerHTML = formatter.format(delta, left);
  }

  _renderResultText(text) {
    let resultPanel = this._resultPanel();
    resultPanel.innerHTML = `<div class="text">${escapeHTML(text)}</div>`;
  }
}

customElements.define('playground-result-panel', PlaygroundResultPanel);
