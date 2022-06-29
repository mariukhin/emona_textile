// node modules
import { action, makeObservable, observable } from 'mobx';
// mocks
import { mockedCatalogItems } from './mocks';

export default class CatalogStore {
  // properties
  @observable catalogItems: CatalogData[] | null = mockedCatalogItems;

  constructor() {
    makeObservable(this);
  }
  
  // actions
  @action setCatalogItems = (data: CatalogData[]) => {
    this.catalogItems = data;
  };
}
