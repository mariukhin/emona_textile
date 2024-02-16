// store
import CatalogItemStore from './store';

class CatalogItemService extends CatalogItemStore {
  getCurrentCatalogItem = (title: string | null) => {
    return this.catalogItemsData.find(item => item.title === title) || null;
  }
}

export default new CatalogItemService();
