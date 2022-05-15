// node modules
import { createBrowserHistory } from 'history';
import { syncHistoryWithStore } from 'mobx-react-router';
import { stringify, StringifyOptions } from 'query-string';
// store
import RoutingStore from './store';

class RoutingService extends RoutingStore {
  public closeWindow = () => {
    window.close();
  };

  public goToRoute = <T extends RouteKeys>(
    link: T,
    options?: {
      params?: BaseObject;
      query?: BaseObject<string, any>;
      inNewTab?: boolean;
    },
  ) => {
    let path = link as unknown as string;

    if (options?.params) {
      path = Object.entries(options.params).reduce(
        (prev, [optionKey, optionValue]) =>
          prev.replace(`:${optionKey}`, optionValue as string),
        path,
      );
    }

    if (options?.query) {
      path = `${path}?${this.stringifyQuery(options.query)}`;
    }

    if (options?.inNewTab) {
      return this.openInNewTab(path);
    }

    this.push(path);
  };

  public goBack = (useHistory = false) => {
    if (useHistory) {
      this.history.goBack();
    } else {
      this.goBack();
    }
  };

  public openInNewTab = (url: string) => window.open(url, '_blank');

  // public changeQuery = <K extends QueryRouterKeys>(
  //   key: keyof QueryRoutes[K],
  //   value: string,
  // ) => {
  //   const search = stringify({
  //     ...this.getQuery(),
  //     [key]: value,
  //   });

  //   this.push({ search });
  // };

  // public getQuery = <K extends QueryRouterKeys>(search?: string) =>
  //   parse(search || this.location.search) as QueryRoutes[K];

  // public deleteQuery<K extends QueryRouterKeys>(key: keyof QueryRoutes[K]) {
  //   const queryObj = this.getQuery<K>();

  //   delete queryObj[key];

  //   const queryString = this.stringifyQuery(queryObj);

  //   this.push({ search: queryString });
  // }

  isOnRoute = (route: RouteKeys) => this.location.pathname.includes(route);

  public stringifyQuery = <Q extends BaseObject>(
    query: Q,
    config?: StringifyOptions,
  ) => {
    const stringifiedQuery = stringify(query, config);

    if (stringifiedQuery) return stringifiedQuery;

    return '';
  };
}

const _RoutingService = new RoutingService();

export const history = syncHistoryWithStore(
  createBrowserHistory(),
  _RoutingService,
);

export default _RoutingService;
