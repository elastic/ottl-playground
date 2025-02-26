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

import {LitElement, html, css, nothing} from 'lit';

class PlaygroundUserConsentBanner extends LitElement {
  static properties = {
    _hasUserConsent: {state: true},
  };

  constructor() {
    super();
    this._hasUserConsent = this._getCookie('userConsent');
  }

  static styles = css`
    .user-consent-banner {
      position: fixed;
      bottom: 0;
      left: 0;
      z-index: 2147483645;
      box-sizing: border-box;
      width: 100%;
      background-color: rgb(245, 245, 245);
      max-height: 600px;
      overflow-y: auto;
      border-top-style: outset;
      border-top-width: 1px;
    }

    .user-consent-banner-inner {
      width: 65%;
      margin: 0 auto;
      padding: 10px;
    }

    .user-consent-banner-copy {
      margin-bottom: 16px;
    }

    .user-consent-banner-header {
      margin-bottom: 8px;
      font-weight: 500;
      font-size: 14px;
      line-height: 15px;
    }

    .user-consent-banner-description {
      font-weight: normal;
      color: #838f93;
      font-size: 12px;
      line-height: 13px;
      text-align: justify;
      width: 100%;
    }

    .user-consent-banner-button {
      box-sizing: border-box;
      display: inline-block;
      min-width: 130px;
      padding: 7px;
      border-radius: 4px;
      background-color: #0073ce;
      border: none;
      cursor: pointer;
      color: #fff;
      text-decoration: none;
      text-align: center;
      font-weight: normal;
      font-size: 13px;
      line-height: 13px;
    }

    .user-consent-banner-button--secondary {
      padding: 9px 13px;
      border: 2px solid #3a4649;
      background-color: transparent;
      color: #0073ce;
    }

    .user-consent-banner-button:hover {
      box-shadow: 0 0 0 999px inset rgba(0, 0, 0, 0.1) !important;
    }

    .user-consent-banner-button--secondary:hover {
      border-color: #838f93;
      background-color: transparent;
      box-shadow: 0 0 0 999px inset rgba(0, 0, 0, 0.1) !important;
    }
  `;

  render() {
    return this._hasUserConsent
      ? nothing
      : html`
          <div class="user-consent-banner">
            <div class="user-consent-banner-inner">
              <div class="user-consent-banner-copy">
                <div class="user-consent-banner-header">NOTICE</div>
                <div class="user-consent-banner-description">
                  <slot name="content">
                    <p>
                      You are solely responsible for any input you submit to
                      this website. We advise that you refrain from submitting
                      any confidential information. If you receive inputs
                      previously submitted by someone else, please make sure to
                      check such content before resubmitting it yourself. The
                      website only provides an execution service, and it does
                      not verify your inputs for quality, security, safety,
                      privacy or any other purpose. Your inputs are only
                      processed for the purpose of providing the execution
                      service you request by pressing the "Run" button, and to
                      allow you to make the actions you perform on the website
                      repeatable and sharable by using the "Copy link" button.
                      Submitting harmful or law-infringing inputs and/or seeking
                      to produce harmful or law-infringing outputs is strictly
                      prohibited, and this website disclaims any liability for
                      such user actions. To maintain the functionality of the
                      service and to monitor its usage for statistical web
                      analytics purposes, this website collects basic usage
                      telemetry data. Such telemetry is limited to network
                      addresses and software agent identifiers. It excludes any
                      user-identifiable personal information.
                    </p>
                    <p>
                      By using this website, you acknowledge the above terms and
                      restrictions.
                    </p>
                  </slot>
                </div>
              </div>
              <div>
                <button
                  class="user-consent-banner-button"
                  @click="${this.acknowledgeNotice}"
                >
                  Acknowledge
                </button>
              </div>
            </div>
          </div>
        `;
  }

  acknowledgeNotice() {
    this._setCookie('userConsent', 'true', 365);
    this._hasUserConsent = true;
  }

  _setCookie(name, value, days) {
    let expires = '';
    if (days) {
      const date = new Date();
      date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
      expires = '; expires=' + date.toUTCString();
    }
    document.cookie = name + '=' + (value || '') + expires + '; path=/';
  }

  _getCookie(name) {
    const nameEQ = name + '=';
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  }
}

customElements.define(
  'playground-user-consent-banner',
  PlaygroundUserConsentBanner
);
