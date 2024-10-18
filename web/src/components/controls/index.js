import {css, html, LitElement} from 'lit-element';
import {globalStyles} from '../../styles.js';
import {repeat} from 'lit/directives/repeat.js';
import {nothing} from 'lit';

export class PlaygroundControls extends LitElement {
  static properties = {
    title: {type: String},
    evaluator: {type: String},
    evaluators: {type: Object},
    hideEvaluators: {type: Boolean, attribute: 'hide-evaluators'},
    hideRunButton: {type: Boolean, attribute: 'hide-run-button'},
    loading: {type: Boolean},
  };

  constructor() {
    super();
    this.title = 'OTTL Playground';
    this.hideEvaluators = false;
    this.hideRunButton = false;
    this.loading = false;
    this.evaluator = 'transform_processor';

    this.evaluators = [
      {id: 'transform_processor', name: 'Transform processor'},
      {id: 'filter_processor', name: 'Filter processor'},
    ];

    window.addEventListener('keydown', (event) => {
      if (event.ctrlKey && event.key.toUpperCase() === 'R') {
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
          width: 70px;
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

        #evaluator {
          width: 160px;
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
            <strong>
              <slot name="app-title-text"><span>${this.title}</span></slot>
            </strong>
          </div>
        </slot>
        <div class="right">
          ${this.hideEvaluators
            ? nothing
            : html`
                <div class="evaluator-container">
                  <label for="evaluator">Evaluator</label>
                </div>
                <div class="evaluator-container">
                  <select
                    id="evaluator"
                    ?disabled="${this.loading}"
                    name="evaluator"
                    @change="${this._notifyEvaluatorChanged}"
                  >
                    ${this.evaluators &&
                    repeat(
                      this.evaluators,
                      (it) => it.id,
                      (it) => {
                        return html` <option
                          ?selected="${it.id === this.evaluator}"
                          value="${it.id}"
                        >
                          ${it.name}
                        </option>`;
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
                              >Run ►
                              <span
                                class="tooltip-text tooltip-text-position-bottom"
                                >⌃+R</span
                              >
                            </span>
                          `}
                    </span>
                  </button>
                </div>
              `}
          <slot name="custom-components"></slot>
        </div>
      </div>
    `;
  }

  _notifyRunRequested() {
    this.dispatchEvent(
      new Event('playground-run-requested', {bubbles: true, composed: true})
    );
  }

  _notifyEvaluatorChanged(e) {
    const event = new CustomEvent('evaluator-changed', {
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
