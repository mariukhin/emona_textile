// node modules
import React from 'react';
import { render } from 'react-dom';
// modules
import { Router } from '../routing/Routing';
// import { Provider } from 'mobx-react';
// import { configure } from 'mobx';

// configure({
//   computedRequiresReaction: true,
//   enforceActions: 'observed',
// });

render(
  // <Provider {...STORES}>
  <Router />,
  // </Provider>,
  window.document.getElementById('root'),
);
