// SPDX-License-Identifier: Apache-2.0

import {css, html, LitElement} from 'lit-element';
import {globalStyles} from '../../styles';

export class PlaygroundCopyLinkButton extends LitElement {
  static properties = {
    label: {},
    loading: {type: Boolean},
    buttonTip: {state: true},
  };

  constructor() {
    super();
    this.label = 'Copy link';
    this.loading = false;
  }

  static get styles() {
    return [
      globalStyles,
      css`
        .link-button {
          background-color: #e8e7e7;
          border: 1px solid #dcdbdb;
          color: black;
          padding: 10px 4px 10px 4px;
          width: 110px;
          text-align: center;
          text-decoration: none;
          font-size: 16px;
          cursor: pointer;
          border-radius: 4px;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .link-button-label {
          color: #4b4949;
        }

        .link-button:hover {
          background-color: #dedede;
        }

        .tooltip:hover .tooltip-text {
          visibility: hidden;
        }
      `,
    ];
  }

  render() {
    return html`
      <button
        id="copy-link-button"
        class="link-button"
        @click="${this._handleCopyLinkClick}"
        ?disabled="${this.loading}"
      >
        <span class="tooltip">
          <span class="link-button-label">${this.label}&nbsp; </span>
          <span
            id="copied-tooltip"
            class="tooltip-text tooltip-text-position-bottom"
            >Copied</span
          >
        </span>
        ${this.loading === true
          ? html`
              <!-- prettier-ignore -->
              <svg height="8px" id="Layer_1" style="enable-background:new 0 0 30 30;" viewBox="0 0 24 20" width="18px" x="0px" xml:space="preserve" xmlns="http://www.w3.org/2000/svg" y="0px"> <rect fill="#000" height="24" width="4" x="0" y="0"> <animate attributeName="opacity" attributeType="XML" begin="0s" dur="0.6s" repeatCount="indefinite" values="1; .2; 1"></animate> </rect> <rect fill="#000" height="24" width="4" x="7" y="0"> <animate attributeName="opacity" attributeType="XML" begin="0.2s" dur="0.6s" repeatCount="indefinite" values="1; .2; 1"></animate> </rect> <rect fill="#000" height="24" width="4" x="14" y="0"> <animate attributeName="opacity" attributeType="XML" begin="0.4s" dur="0.6s" repeatCount="indefinite" values="1; .2; 1"></animate> </rect> </svg>
            `
          : html`
              <!-- prettier-ignore -->
              <svg width="17px" height="17px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M13.5442 10.4558C11.8385 8.75022 9.07316 8.75022 7.36753 10.4558L4.27922 13.5442C2.57359 15.2498 2.57359 18.0152 4.27922 19.7208C5.98485 21.4264 8.75021 21.4264 10.4558 19.7208L12 18.1766" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> <path d="M10.4558 13.5442C12.1614 15.2498 14.9268 15.2498 16.6324 13.5442L19.7207 10.4558C21.4264 8.75021 21.4264 5.98485 19.7207 4.27922C18.0151 2.57359 15.2497 2.57359 13.5441 4.27922L12 5.82338" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> </g></svg>
            `}
      </button>
    `;
  }

  async _handleCopyLinkClick() {
    this.loading = true;
    setTimeout(() => {
      let copied = this.dispatchEvent(
        new CustomEvent('copy-link-click', {
          composed: true,
          cancelable: true,
        })
      );

      if (copied) {
        this.shadowRoot.querySelector('#copied-tooltip').style.visibility =
          'visible';
        setTimeout(() => {
          this.shadowRoot.querySelector('#copied-tooltip').style.visibility =
            'hidden';
        }, 1500);
      }
      this.loading = false;
    }, 0);
  }
}

customElements.define('playground-copy-link-button', PlaygroundCopyLinkButton);
