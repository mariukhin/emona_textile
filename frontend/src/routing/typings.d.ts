import { ROUTES } from './registration';

declare global {
  type RouteKeys = valueof<typeof ROUTES>;
}
