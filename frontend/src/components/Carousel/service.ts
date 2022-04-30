// modules
// store
import CarouselStore from './store';

class CarouselService extends CarouselStore {
  changeCurrentItem = (direction: Directions) => {
    if (direction === 'Right') {
      let updatedItemIdx = -1;
      const mappedData = this.carouselItems?.map((item, idx) => {
        if (item.isCurrent) {
          updatedItemIdx = idx + 1;
  
          return {
            ...item,
            isCurrent: false,
          }
        }

        if (idx === updatedItemIdx) {
          return {
            ...item,
            isCurrent: true,
          }
        }

        return item;
      });

      console.log(mappedData);
      
      this.setCarouselItems(mappedData || []);
    }
  }
}

export default new CarouselService();
