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
        padding: 10px 3px 10px 10px;
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
              <svg viewBox="0 0 400 360" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg" version="1.1" xml:space="preserve"><g class="layer"><g id="e2352819-c193-45a2-976b-d3c08802fe29" transform="translate(-5.44 -5.93) translate(-59 -37) translate(-252 -247) matrix(1 0 0 1 540 540)"/> <g id="SVGRepo_iconCarrier" transform="translate(-5.44 -5.93) translate(-59 -37) translate(-252 -247) matrix(0.76 -0.32 0.32 0.76 487.17 478.02)"> <path d="m382.4,75.14l-259.47,54.02c-4.79,1 -8.99,3.85 -11.67,7.95c-2.69,4.09 -3.63,9.08 -2.63,13.87l0.02,0.12l-62.78,13.08l17.13,82.27l62.78,-13.07l0.03,0.13c1.81,8.71 9.48,14.69 18.04,14.69c1.25,0 2.51,-0.13 3.78,-0.39l65.85,-13.71c-1.54,-5.82 -2.37,-11.92 -2.37,-18.22c0,-39.27 31.95,-71.22 71.22,-71.22c35.62,0 65.21,26.27 70.42,60.45l54.36,-11.32l-24.71,-118.65z" fill="#ffffff" fill-rule="nonzero" id="svg_7" stroke-dashoffset="0" stroke-miterlimit="4" transform=" translate(-226.49, -161.67)"/> </g> <g id="svg_1" transform="translate(-5.44 -5.93) translate(-59 -37) translate(-252 -247) matrix(0.76 -0.32 0.32 0.76 547.01 501.09)"> <path d="m282.33,157.79c-32.03,0 -58.09,26.06 -58.09,58.09c0,32.04 26.06,58.1 58.09,58.1c32.04,0 58.1,-26.06 58.1,-58.1c0,-32.03 -26.06,-58.09 -58.1,-58.09zm0,79.3c-11.69,0 -21.2,-9.51 -21.2,-21.2c0,-11.7 9.51,-21.21 21.2,-21.21c11.69,0 21.21,9.51 21.21,21.21c0,11.68 -9.52,21.2 -21.21,21.2z" fill="#ffffff" fill-rule="nonzero" id="svg_8" stroke-dashoffset="0" stroke-miterlimit="4" transform=" translate(-282.33, -215.88)"/> </g> <g id="svg_2" transform="translate(-5.44 -5.93) translate(-59 -37) translate(-252 -247) matrix(0.76 -0.32 0.32 0.76 645.6 375.48)"> <path d="m511.61,175.89l-26.24,-126.05c-1,-4.79 -3.86,-8.98 -7.95,-11.66c-4.09,-2.69 -9.08,-3.63 -13.87,-2.64l-56.58,11.78c-3.13,0.65 -5.9,2.06 -8.17,3.99c-4.98,4.23 -7.56,10.98 -6.13,17.83l26.24,126.05c1,4.79 3.86,8.99 7.95,11.67c1.49,0.98 3.11,1.72 4.79,2.23c1.72,0.52 3.52,0.79 5.32,0.79c1.25,0 2.51,-0.13 3.76,-0.39l56.58,-11.78c9.97,-2.08 16.38,-11.85 14.3,-21.82z" fill="#ffffff" fill-rule="nonzero" id="svg_9" stroke-dashoffset="0" stroke-miterlimit="4" transform=" translate(-452.14, -122.52)"/> </g> <g id="svg_3" transform="translate(-5.44 -5.93) translate(-59 -37) translate(-252 -247) matrix(0.76 -0.32 0.32 0.76 350.11 581.08)"> <path d="m48.48,236.55l-11.97,-57.5c-1.65,-7.9 -8.12,-13.54 -15.69,-14.52c-2,-0.26 -4.06,-0.21 -6.13,0.22c-9.97,2.07 -16.38,11.84 -14.3,21.82l10.79,51.83l1.18,5.67c1.81,8.71 9.49,14.69 18.04,14.69c1.25,0 2.51,-0.13 3.78,-0.39c2.08,-0.43 3.98,-1.21 5.71,-2.24c6.55,-3.92 10.24,-11.68 8.59,-19.58z" fill="#ffffff" fill-rule="nonzero" id="svg_10" stroke-dashoffset="0" stroke-miterlimit="4" transform=" translate(-24.43, -211.57)"/> </g> <g id="svg_4" transform="translate(-5.44 -5.93) translate(-59 -37) translate(-252 -247) matrix(0.7 0 0 0.56 548.23 589.67)"> <path d="m378.59,452.58l-57.47,-176.98c-10.2,6.65 -22.19,11.06 -35.08,11.72l31.64,97.73l-70.69,0l31.64,-97.74c-12.89,-0.66 -24.88,-4.93 -35.08,-11.58l-57.46,176.89c-3.15,9.69 2.15,20.06 11.84,23.2c9.69,3.15 20.1,-1.88 23.25,-11.57l13.83,-42.31l94.66,0l13.83,42.3c2.53,7.79 9.76,12.6 17.54,12.6c1.88,0 3.81,-0.36 5.7,-0.98c9.69,-3.14 14.99,-13.59 11.85,-23.28z" fill="#ffffff" fill-rule="nonzero" id="svg_11" stroke-dashoffset="0" stroke-miterlimit="4" transform=" translate(-282.34, -376.22)"/> </g> </g></svg>
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
