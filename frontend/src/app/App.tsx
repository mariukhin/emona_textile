// node modules
import React from 'react';
import { createRoot } from 'react-dom/client';
// modules
import { Router } from '../routing/Routing';
import { Provider } from 'mobx-react';
import { configure } from 'mobx';
import STORES from './stores';

configure({
  computedRequiresReaction: true,
  enforceActions: 'observed',
});

const container = document.getElementById('root');
const root = createRoot(container!);

root.render(
  <Provider {...STORES}>
    <Router />,
  </Provider>,
);
