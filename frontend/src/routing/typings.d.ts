import { ROUTES } from './registration';

declare global {
  type RouteKeys = valueof<typeof ROUTES>;

  type QueryRoutes = {
    BASE: {
      theme?: ThemeKeys;
    };
    CHILD_CHART: {
      exchangeToken: string;
      interval: TradingChartUserPeriod;
      library: TradingChartLibraries;
    };
    STANDALONE_CHART: {
      exchange: string;
      token: string;
      chartType: ChartIQType
    }
    AUTH: {
      code: string;
    };
    TRADE_ORDER: Omit<TradeButtonBaseOrder, 'amo'> & {
      amo?: 'true';
      orders?: string;
    };
  };
  
  type QueryRouterKeys = keyof QueryRoutes;
}
