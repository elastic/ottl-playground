// SPDX-License-Identifier: Apache-2.0

import {css} from 'lit-element';

export const globalStyles = css`
  input[type='text'],
  select,
  textarea {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
    resize: vertical;
  }

  button:disabled,
  button[disabled] {
    opacity: 0.6;
    cursor: not-allowed;
  }

  label {
    padding: 10px 10px 10px 0;
    display: inline-block;
  }

  .hidden-overflow {
    overflow: hidden;
  }

  .h-full {
    height: 100%;
  }

  .w-full {
    width: 100%;
  }

  .full-size {
    height: 100%;
    width: 100%;
  }

  .tooltip {
    position: relative;
    display: inline-block;
  }

  .tooltip .tooltip-text {
    visibility: hidden;
    width: 100%;
    background-color: rgb(84, 84, 84);
    color: #fff;
    text-align: center;
    padding: 5px 0;
    border-radius: 6px;
    position: absolute;
    z-index: 1;
    font-size: small;
  }

  .tooltip:hover .tooltip-text {
    visibility: visible;
  }

  .tooltip-text-position-right {
    top: -5px;
    left: 105%;
  }

  .tooltip-text-position-left {
    top: -5px;
    right: 105%;
  }

  .tooltip-text-position-bottom {
    width: 120px;
    top: 100%;
    left: 50%;
    margin-left: -60px;
  }
`;
