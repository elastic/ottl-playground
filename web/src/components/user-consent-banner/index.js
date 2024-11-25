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
      background-color: #f1f6f4;
      max-height: 600px;
      overflow-y: auto;
      border-top-style: outset;
      border-top-width: 1px;
    }

    .user-consent-banner-inner {
      width: 65%;
      margin: 0 auto;
      padding: 32px;
    }

    .user-consent-banner-copy {
      margin-bottom: 16px;
    }

    .user-consent-banner-actions {
    }

    .user-consent-banner-header {
      margin-bottom: 8px;
      font-weight: bold;
      font-size: 16px;
      line-height: 24px;
    }

    .user-consent-banner-description {
      font-weight: normal;
      color: #838f93;
      font-size: 16px;
      line-height: 19px;
      text-align: justify;
      width: 100%;
    }

    .user-consent-banner-button {
      box-sizing: border-box;
      display: inline-block;
      min-width: 164px;
      padding: 11px 13px;
      border-radius: 4px;
      background-color: #0073ce;
      border: none;
      cursor: pointer;
      color: #fff;
      text-decoration: none;
      text-align: center;
      font-weight: normal;
      font-size: 16px;
      line-height: 20px;
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
                      privacy or any other purpose.
                    </p>
                    <p>
                      Your inputs are only processed for the purpose of
                      providing the execution service you request by pressing
                      the "Run" button, and to allow you to make the actions you
                      perform on the website repeatable and sharable by using
                      the "Copy link" button.
                    </p>
                    <p>
                      Submitting harmful or law-infringing inputs and/or seeking
                      to produce harmful or law-infringing outputs is strictly
                      prohibited, and this website disclaims any liability for
                      such user actions.
                    </p>
                    <p>
                      To maintain the functionality of the service and to
                      monitor its usage for statistical web analytics purposes,
                      this website collects basic usage telemetry data. Such
                      telemetry is limited to network addresses and software
                      agent identifiers. It excludes any user-identifiable
                      personal information.
                    </p>
                    By using this website, you acknowledge the above terms and
                    restrictions.
                  </slot>
                </div>
              </div>
              <div class="user-consent-banner-actions">
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
