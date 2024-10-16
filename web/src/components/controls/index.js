import {css, html, LitElement} from 'lit-element';
import {globalStyles} from '../../styles.js';
import {nothing} from 'lit';

export class PlaygroundControls extends LitElement {
  static properties = {
    // attributes
    title: {type: String, attribute: true},
    evaluator: {type: String, attribute: true, reflect: true},
    hideEvaluator: {type: Boolean, attribute: 'hide-evaluator'},
    hideRunButton: {type: Boolean, attribute: 'hide-run-button'},
    // states
    _loading: {type: Boolean, state: true},
  };

  constructor() {
    super();
    this.title = 'OTTL Playground';
    this._loading = true;
    this._addEventListeners();
  }

  static get styles() {
    return [
      globalStyles,
      css`
        .playground-controls {
          overflow: hidden;
          width: 100%;
        }

        .playground-controls div {
          float: left;
          text-align: center;
          margin: 10px 0 10px 0;
          text-decoration: none;
          font-size: 17px;
        }

        .playground-controls .app-title {
          font-size: 25px;
          padding: 15px;
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
          border: none;
          color: white;
          padding: 10px;
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
      `,
    ];
  }

  get selectedEvaluator() {
    return this.shadowRoot.querySelector('#evaluator').value;
  }

  render() {
    return html`
      <style>
        ${this.hideEvaluator === true
          ? '.evaluator-container { visibility: hidden; }'
          : nothing}
      </style>
      <div class="playground-controls">
        <slot name="app-title">
          <div class="app-title">
            <strong>
              <slot name="app-title-text"><span>${this.title}</span></slot>
            </strong>
          </div>
        </slot>
        <div class="right">
          ${this.hideEvaluator !== true
            ? html`
                <div class="evaluator-container">
                  <label for="evaluator">Evaluator</label>
                </div>
                <div class="evaluator-container">
                  <select
                    id="evaluator"
                    ?disabled="${this._loading}"
                    name="evaluator"
                    .value=${this.evaluator}
                    @change="${this._handleEvaluatorChange}"
                  >
                    <option selected value="transform_processor">
                      Transform processor
                    </option>
                    <option value="filter_processor">Filter processor</option>
                  </select>
                </div>
              `
            : nothing}
          ${this.hideRunButton !== true
            ? html`
                <div>
                  <button
                    class="run-button"
                    ?disabled="${this._loading}"
                    id="btn-run"
                    @click="${this._notifyRunClick}"
                  >
                    <span id="run-btn">
                      ${this._loading
                        ? html`
                            <span>
                              <svg
                                height="8px"
                                id="Layer_1"
                                style="enable-background:new 0 0 30 30;"
                                viewBox="0 0 18 10"
                                width="10px"
                                x="0px"
                                xml:space="preserve"
                                xmlns="http://www.w3.org/2000/svg"
                                y="0px"
                              >
                                <rect
                                  fill="#fff"
                                  height="20"
                                  width="4"
                                  x="0"
                                  y="0"
                                >
                                  <animate
                                    attributeName="opacity"
                                    attributeType="XML"
                                    begin="0s"
                                    dur="0.6s"
                                    repeatCount="indefinite"
                                    values="1; .2; 1"
                                  ></animate>
                                </rect>
                                <rect
                                  fill="#fff"
                                  height="20"
                                  width="4"
                                  x="7"
                                  y="0"
                                >
                                  <animate
                                    attributeName="opacity"
                                    attributeType="XML"
                                    begin="0.2s"
                                    dur="0.6s"
                                    repeatCount="indefinite"
                                    values="1; .2; 1"
                                  ></animate>
                                </rect>
                                <rect
                                  fill="#fff"
                                  height="20"
                                  width="4"
                                  x="14"
                                  y="0"
                                >
                                  <animate
                                    attributeName="opacity"
                                    attributeType="XML"
                                    begin="0.4s"
                                    dur="0.6s"
                                    repeatCount="indefinite"
                                    values="1; .2; 1"
                                  ></animate>
                                </rect>
                              </svg>
                            </span>
                          `
                        : 'Run ►'}
                    </span>
                  </button>
                </div>
              `
            : nothing}
          <slot name="custom-components"></slot>
        </div>
      </div>
    `;
  }

  _notifyRunClick() {
    this.dispatchEvent(
      new Event('playground-run', {bubbles: true, composed: true})
    );
  }

  _addEventListeners() {
    let that = this;
    window.addEventListener('playground-wasm-ready', function () {
      that._loading = false;
    });

    window.addEventListener('playground-evaluator-change', (e) => {
      this.evaluator = e.detail.value;
    });
  }

  _handleEvaluatorChange(e) {
    const event = new CustomEvent('playground-evaluator-change', {
      detail: {value: e.target.value},
      bubbles: true,
      composed: true,
      cancelable: true,
    });
    this.dispatchEvent(event);
  }
}

customElements.define('playground-controls', PlaygroundControls);
