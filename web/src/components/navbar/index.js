// SPDX-License-Identifier: Apache-2.0

import {css, html, LitElement} from 'lit-element';
import {nothing} from 'lit';

export class PlaygroundNavBar extends LitElement {
  static properties = {
    title: {type: String},
    githubLink: {type: String, attribute: 'github-link'},
  };

  constructor() {
    super();
    this.title = 'OTTL&nbsp;Playground';
  }

  static get styles() {
    return css`
      :host .navbar {
        width: 100%;
        overflow: auto;
        display: flex;
        background-color: var(--background-color, rgb(79, 98, 173));
        border-bottom: 2px solid var(--border-bottom-color, #333f70);
      }

      .navbar a {
        text-align: center;
        padding: 10px 5px 10px 10px;
        color: white;
        text-decoration: none;
        font-size: 17px;
      }

      .navbar a .logo svg {
        height: 30px;
      }

      @media screen and (max-width: 500px) {
        .navbar a {
          float: none;
          display: block;
        }
      }

      .github-corner:hover .octo-arm {
        animation: octocat-wave 560ms ease-in-out;
      }

      @keyframes octocat-wave {
        0%,
        100% {
          transform: rotate(0);
        }
        20%,
        60% {
          transform: rotate(-25deg);
        }
        40%,
        80% {
          transform: rotate(10deg);
        }
      }

      @media (max-width: 500px) {
        .github-corner:hover .octo-arm {
          animation: none;
        }
        .github-corner .octo-arm {
          animation: octocat-wave 560ms ease-in-out;
        }
      }

      .title-beta-box {
        font-size: 7px !important;
        font-weight: 300 !important;
        color: white;
        border: gray solid 1px;
        padding-left: 2px;
        padding-right: 2px;
        margin-top: -8px;
        cursor: default;
      }

      .title {
        display: flex;
        color: white;
        align-items: center;
        gap: 3px;
        width: 100%;
      }

      .title .text {
        font-weight: 500;
        cursor: default;
      }

      .github-link-container {
        height: 100%;
        width: 100%;
        display: flex;
        justify-content: flex-end;
      }
    `;
  }

  render() {
    return html`
      <div class="navbar">
        <a href="">
          <slot name="logo">
            <span class="logo">
              <!-- prettier-ignore -->
              <svg role="img" fill="#fff" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><title>OpenTelemetry</title><path d="M12.6974 13.1173c-1.0224 1.0224-1.0224 2.68 0 3.7024 1.0224 1.0224 2.68 1.0224 3.7024 0 1.0224-1.0223 1.0224-2.68 0-3.7024-1.0223-1.0223-2.68-1.0223-3.7024 0zm2.7677 2.7701c-.5063.5063-1.3267.5063-1.833 0s-.5063-1.3266 0-1.833c.5063-.5062 1.3267-.5062 1.833 0 .5063.504.5063 1.3267 0 1.833zM16.356.2355l-1.6041 1.6042c-.314.314-.314.83 0 1.144L21.015 9.247c.314.314.83.314 1.144 0l1.6042-1.6041c.314-.314.314-.83 0-1.144L17.4976.2354c-.314-.314-.8276-.314-1.1416 0zM5.1173 20.734c.2848-.2848.2848-.7497 0-1.0345l-.8155-.8155c-.2848-.2848-.7497-.2848-1.0345 0l-1.6845 1.6845-.0024.0024-.4625-.4625c-.2556-.2556-.6718-.2556-.925 0-.2556.2556-.2556.6718 0 .925l2.775 2.775c.2556.2556.6718.2556.925 0 .2532-.2556.2556-.6718 0-.925l-.4625-.4625.0024-.0024zm8.4856-15.893-3.5637 3.5637c-.3164.3164-.3164.8374 0 1.1538l2.2006 2.2005c1.5554-1.1197 3.7365-.981 5.1361.4187l1.7819-1.7818c.3164-.3165.3164-.8374 0-1.1538l-4.401-4.401c-.3165-.319-.8374-.319-1.1539 0zm-2.2881 7.8455-1.2999-1.2999c-.3043-.3043-.8033-.3043-1.1076 0l-4.5836 4.586c-.3042.3043-.3042.8033 0 1.1076l2.5973 2.5973c.3043.3043.8033.3043 1.1076 0l2.9478-2.9527c-.6231-1.2877-.5112-2.8431.3384-4.0383z"/></svg>
            </span>
          </slot>
        </a>

        ${this.title
          ? html`
              <div class="title">
                <span class="text">OTTL&nbsp;Playground</span>
                <sup
                  class="title-beta-box"
                  title="The OTTL Playground is still in beta and the authors of this tool would welcome your feedback"
                  >BETA</sup
                >
              </div>
            `
          : nothing}

        <slot name="custom-components"></slot>

        ${!this.githubLink
          ? nothing
          : html`
              <div class="github-link-container">
                <a
                  href="${this.githubLink}"
                  class="github-corner"
                  target="_blank"
                >
                  <!-- prettier-ignore -->
                  <svg width="60" height="60" color="#fff" fill="#333f70" style="border:0;position:absolute;right:0;top:0;" viewBox="0 20 230 230" xmlns="http://www.w3.org/2000/svg">
                        <path d="m0 0 115 115h15l12 27 108 108v-250z"/>
                        <path class="octo-arm" d="m128.3 109c-14.5-9.3-9.3-19.4-9.3-19.4 3-6.9 1.5-11 1.5-11-1.3-6.6 2.9-2.3 2.9-2.3 3.9 4.6 2.1 11 2.1 11-2.6 10.3 5.1 14.6 8.9 15.9" fill="currentColor" style="transform-origin:130px 106px"/>
                        <path class="octo-body" d="m115 115c-0.1 0.1 3.7 1.5 4.8 0.4l13.9-13.8c3.2-2.4 6.2-3.2 8.5-3-8.4-10.6-14.7-24.2 1.6-40.6 4.7-4.6 10.2-6.8 15.9-7 0.6-1.6 3.5-7.4 11.7-10.9 0 0 4.7 2.4 7.4 16.1 4.3 2.4 8.4 5.6 12.1 9.2 3.6 3.6 6.8 7.8 9.2 12.2 13.7 2.6 16.2 7.3 16.2 7.3-3.6 8.2-9.4 11.1-10.9 11.7-0.3 5.8-2.4 11.2-7.1 15.9-16.4 16.4-30 10-40.6 1.6 0.2 2.8-1 6.8-5 10.8l-11.7 11.6c-1.2 1.2 0.6 5.4 0.8 5.3z" fill="currentColor"/>
                      </svg>
                </a>
              </div>
            `}
      </div>
    `;
  }
}

customElements.define('playground-navbar', PlaygroundNavBar);
