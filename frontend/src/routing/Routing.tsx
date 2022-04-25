// node modules
import React, { memo } from 'react';
import { Route, Routes, BrowserRouter } from 'react-router-dom';
// components
import AppLayout from 'app/AppLayout';
// routes
import { ROUTES } from './registration';
import Suspense from './Suspense';

// Lazy Components -_>
// pages
const HomeScreen = React.lazy(() => import('../pages/HomeScreen'));
// const ScripDetailsScreen = React.lazy(() => import('app/pages/ScripDetailsScreen'));
// const UserScreen = React.lazy(() => import('app/pages/UserScreen'));
// const TradeDetailsScreen = React.lazy(() => import('app/pages/TradeDetailsScreen'));
// const OptionChainScreen = React.lazy(() => import('app/pages/OptionChainScreen'));
// const DiscoverScreen = React.lazy(() => import('app/pages/DiscoverScreen'));
// const InboxScreen = React.lazy(() => import('app/pages/InboxScreen'));
// const DiscoverDetails = React.lazy(() => import('app/pages/DiscoverDetails'));
// const NotificationSettingsScreen = React.lazy(() => import('app/pages/NotificationSettingsScreen'));
// const IPODetailsScreen = React.lazy(() => import('app/pages/IPODetailsScreen'));
// const ImportantAlertsScreen = React.lazy(() => import('app/pages/ImportantAlertsScreen'));
// // charts
// const ChildChart = React.lazy(() => import('modules/TradingChart/components/ChildChart'));
// <-- Lazy Components

export const Router = memo(() => (
  <AppLayout>
    <Suspense>
      <BrowserRouter>
        <Routes>
          <Route path={ROUTES.HOME} element={<HomeScreen />} />
          {/* <Route
              exact
              path={AUTHORIZED_ROUTES.PERFORMANCE}
              component={PerformanceScreen}
            />
            <Route
              exact
              path={[AUTHORIZED_ROUTES.HOME, AUTHORIZED_ROUTES.DISCOVER]}
              component={DiscoverScreen}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.DISCOVER_DETAILS}
              component={DiscoverDetails}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.SCRIP_DETAILS}
              component={ScripDetailsScreen}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.USER_WALLET_WITHDRAW}
              component={WithdrawFunds}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.USER_REWARDS_WITHDRAW_CONFIRM}
              component={WithdrawRewardsConfirm}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.USER_REWARDS_WITHDRAW}
              component={WithdrawRewards}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.USER_WALLET_ADD_FUNDS}
              component={AddFunds}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.USER_WALLET_WITHDRAW_CONFIRM}
              component={WithdrawConfirm}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.USER_WALLET_ADD_FUNDS_CONFIRM}
              component={AddFundsMethods}
            />
            <Route
              path={AUTHORIZED_ROUTES.WALLET_DETAILS}
              component={WalletDetails}
            />
            <Route
              path={AUTHORIZED_ROUTES.USER_REWARDS_WALLET}
              component={RewardsWalletDetails}
            />
            <Route path={AUTHORIZED_ROUTES.USER} component={UserScreen} />
            <Route
              path={AUTHORIZED_ROUTES.TRADES_DETAILS}
              component={TradeDetailsScreen}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.OPTION_CHAIN}
              component={OptionChainScreen}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.INBOX}
              component={InboxScreen}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.NOTIFICATION_SETTINGS}
              component={NotificationSettingsScreen}
            />
            <Route
              path={AUTHORIZED_ROUTES.IPO_DETAILS}
              component={IPODetailsScreen}
            />
            <Route
              exact
              path={AUTHORIZED_ROUTES.IMPORTANT_ALERTS}
              component={ImportantAlertsScreen}
            /> */}
        </Routes>
      </BrowserRouter>
    </Suspense>
  </AppLayout>
));
