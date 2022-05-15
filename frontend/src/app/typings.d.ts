import STORES from './stores';

declare global {
  type StoresType = typeof STORES;

   // base object type
  type BaseObject<K extends string = string, V = unknown> = {
    [key in K]: V;
  };
}
