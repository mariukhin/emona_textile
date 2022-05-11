// node modules
import { createBrowserHistory } from 'history';
import { syncHistoryWithStore } from 'mobx-react-router';
// modules
// import { EXTERNAL_LINKS } from 'app/routing/registration';
// import AuthStore from 'modules/Auth';
// store
import RoutingStore from './store';

class RoutingService extends RoutingStore {
  public closeWindow = () => {
    window.close();
  };

  // public goToRoute = <T extends RouteKeys>(
  //   link: T,
  //   options?: {
  //     params?: BaseObject;
  //     query?: BaseObject<string, any>;
  //     inNewTab?: boolean;
  //   },
  // ) => {
  //   let path = link as string;

  //   if (options?.params) {
  //     path = Object.entries(options.params).reduce(
  //       (prev, [optionKey, optionValue]) =>
  //         prev.replace(`:${optionKey}`, optionValue as string),
  //       path,
  //     );
  //   }

  //   if (options?.query) {
  //     path = `${path}?${this.stringifyQuery(options.query)}`;
  //   }

  //   if (options?.inNewTab) {
  //     return this.openInNewTab(path);
  //   }

  //   this.push(path);
  // };

  // goToExternalRoute = <Q extends BaseObject>(
  //   linkKey: ExternalRoutesKeys,
  //   config: ExternalRoutesOptions<Q> = {},
  // ) => {
  //   const { query, target, withToken, betaStyle } = config;

  //   const fullUrl = this.getExternalLinkByKey(linkKey);

  //   const [path, hash] = fullUrl.split('#');

  //   let url = path;

  //   if (query || withToken || betaStyle) {
  //     url += '?';

  //     const _query: BaseObject = { ...(query || {}) };

  //     if (withToken) {
  //       _query.at = AuthStore.loginKey;
  //     }

  //     if (betaStyle) {
  //       _query.appvariant = '4.0';
  //     }

  //     url += this.stringifyQuery(_query);
  //   }

  //   if (hash) {
  //     url += `#${hash}`;
  //   }

  //   return window.open(url, target);
  // };

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

  // public stringifyQuery = <Q extends BaseObject>(
  //   query: Q,
  //   config?: StringifyOptions,
  // ) => {
  //   const stringifiedQuery = stringify(query, config);

  //   if (stringifiedQuery) return stringifiedQuery;

  //   return '';
  // };
}

const _RoutingService = new RoutingService();

export const history = syncHistoryWithStore(
  createBrowserHistory(),
  _RoutingService,
);

export default _RoutingService;
