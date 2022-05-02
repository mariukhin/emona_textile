// node modules
import * as R from 'ramda';
import { action, computed, makeObservable, observable } from 'mobx';
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

  // computed
  @computed get currentItem() {
    return this.carouselItems && R.find((item) => item.isCurrent, this.carouselItems);
  }
}
