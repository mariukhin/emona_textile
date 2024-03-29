// node modules
import React from 'react';
import ReactGA from 'react-ga';
import { createRoot } from 'react-dom/client';
// modules
import { CustomRouter as Router } from '../routing/Routing';
import { Provider } from 'mobx-react';
import { configure } from 'mobx';
import STORES from './stores';
import 'axios/axiosDefaults';
import './style.css';

const TRACKING_ID = "G-KL3H0GMWKB";
ReactGA.initialize(TRACKING_ID);

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
