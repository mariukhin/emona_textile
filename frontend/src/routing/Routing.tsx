// node modules
import React, { memo } from 'react';
import { Router, Route, Switch } from 'react-router-dom';
// components
import AppLayout from 'app/AppLayout';
// routes
import { ROUTES } from './registration';
import Suspense from './Suspense';
import { history } from './service';

// Lazy Components -_>
// pages
const HomeScreen = React.lazy(() => import('../pages/HomeScreen'));
// <-- Lazy Components

export const CustomRouter = memo(() => (
  <Router history={history}>
    <AppLayout>
      <Suspense>
        <Switch>
          <Route exact path={ROUTES.HOME} component={HomeScreen} />
          {/* <Route
              exact
              path={AUTHORIZED_ROUTES.PERFORMANCE}
              component={PerformanceScreen}
            />
           */}
        </Switch>
      </Suspense>
    </AppLayout>
  </Router>
));
