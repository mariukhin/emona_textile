// modules
import axios from 'axios';
// store
import CatalogStore from './store';

class CatalogService extends CatalogStore {
  getCatalogItems = async () => {
    const { data } = await axios.get('/catalog');

    if (data) this.setCatalogItems(data);
  }
}

export default new CatalogService();
