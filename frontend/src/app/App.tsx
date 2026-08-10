// node modules
import React from 'react';
import ReactGA from 'react-ga';
import { createRoot } from 'react-dom/client';
// modules
import { CustomRouter as Router } from '../routing/Routing';
import { Provider } from 'mobx-react';
import { configure } from 'mobx';
import STORES from './stores';
import { GA_TRACKING_ID } from 'config';
import 'axios/axiosDefaults';
import './style.css';

if (GA_TRACKING_ID) {
  ReactGA.initialize(GA_TRACKING_ID);
}

configure({
  computedRequiresReaction: true,
  enforceActions: 'observed',
});

const container = document.getElementById('root');
const root = createRoot(container!);

root.render(
  <Provider {...STORES}>
    <Router />
  </Provider>,
);
