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
`;
