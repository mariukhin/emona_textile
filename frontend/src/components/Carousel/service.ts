// modules
import * as R from 'ramda';
import axios from 'axios';
// store
import CarouselStore from './store';

class CarouselService extends CarouselStore {
  getCarouselItems = async () => {
    const { data } = await axios.get('/carousel');

    if (data) {
      this.setCarouselItems(data);
    } 
  }

  changeCurrentItem = (direction: Directions) => {
    if (!this.carouselItems) return null;

    let updatedItemIdx = -1;
    const itemsLength = this.carouselItems.length;

    const result = this.carouselItems.map((item, idx) => {
      if (item.isCurrent) {
        updatedItemIdx = direction === 'Right' ? idx + 1 : idx - 1;

        if (updatedItemIdx === itemsLength) {
          updatedItemIdx = 0;
        }

        if (updatedItemIdx == -1) {
          updatedItemIdx = itemsLength - 1;
        }

        return {
          ...item,
          isCurrent: false,
        };
      }

      return item;
    });

    const updatedCarouselItems: CarouselData[] =
      R.update(updatedItemIdx, { ...result[updatedItemIdx], isCurrent: true }, result);

    this.setCarouselItems(updatedCarouselItems);
  };
}

export default new CarouselService();
