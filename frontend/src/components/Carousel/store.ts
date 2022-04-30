// node modules
import { action } from 'mobx';
// mocks
import { mockedCarouselItems } from './mocks';

export default class CarouselStore {
  // properties
  carouselItems: CarouselData[] | null = mockedCarouselItems;
  
  
  // actions
  @action setCarouselItems = (data: CarouselData[]) => {
    this.carouselItems = data;
    console.log('sss', this.carouselItems);
  };
}
