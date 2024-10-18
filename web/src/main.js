import './components/navbar/index.js';
import './components/controls/index.js';
import './components/playground.js';

window.addEventListener('load', () => {
  let button = document.querySelector('#myButton');
  button.addEventListener('click', () => {
    let main = document.querySelector('playground-stage');
    let a = {...main.state};
    a.evaluator = 'filter_processor';
    a.payloadType = 'metrics';
    a.config = '-- this is my config --';
    a.payload = '{ "s": 1 }';
    main.state = a;
  });
});
