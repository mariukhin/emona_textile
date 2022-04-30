import STORES from './stores';

declare global {
  type StoresType = typeof STORES;
}
