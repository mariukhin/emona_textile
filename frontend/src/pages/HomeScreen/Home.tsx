// node modules
import React from 'react';
import { observer } from 'mobx-react';
// components
import Carousel from 'components/Carousel';

const HomePageView = () => (
  <Carousel />
);

export default observer(HomePageView);
