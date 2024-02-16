// node modules
import { action, makeObservable, observable } from 'mobx';
// mocks
import { mockedCatalogItems } from './mocks';

export default class CatalogItemStore {
  // properties
  @observable catalogItemsData: CatalogItemData[] = mockedCatalogItems;
  @observable currentCatalogItem: CatalogItemData | null = null;

  constructor() {
    makeObservable(this);
  }
  
  // actions
  @action setCurrentCatalogItem = (title: CatalogData['title']) => {
    this.currentCatalogItem = this.catalogItemsData.find(item => item.title === title) || null;
  };
}
