// node modules
import { action, makeObservable, observable } from 'mobx';
// mocks
import { mockedCarouselItems } from './mocks';

export default class CarouselStore {
  // properties
  @observable carouselItems: CarouselData[] | null = mockedCarouselItems;

  constructor() {
    makeObservable(this);
  }
  
  // actions
  @action setCarouselItems = (data: CarouselData[]) => {
    this.carouselItems = data;
  };
}
