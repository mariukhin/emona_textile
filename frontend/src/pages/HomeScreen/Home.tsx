// node modules
import React, { useState } from 'react';
import Ticker from 'react-ticker';
import * as R from 'ramda';
// components
import AboutUsBlock from 'components/AboutUsBlock';
import AdvantagesBlock from 'components/AdvantagesBlock';
import BlockInfoComponent from 'components/BlockInfoComponent';
import { Carousel } from 'components/Carousel';
import { Catalog } from 'components/Catalog';
import { ContactsAndFormBlock } from 'components/ContactsAndFormBlock';
// mocks
import { mockedAdvantages } from 'components/AdvantagesBlock/mocks';
import { logos } from './mocks';
// styles
import {
  ClientsBlockWrapper,
  ClientsBlockTitleWrapper,
  ClientTicker,
  TickerWrapper,
  TickerImage,
} from './styles';

const HomePageView = () => {
  const [isMouseOverTicker, setIsMouseOverTicker] = useState(true);

  const handleMouseOver = () => {
    setIsMouseOverTicker(false);
  }

  const handleMouseLeave = () => {
    setIsMouseOverTicker(true);
  }

  const shuffler = R.curry(function(random, list) {
    let idx = -1;
    let len = list.length;
    let position;
    let result: any[] = [];
    while (++idx < len) {
        position = Math.floor((idx + 1) * random());
        result[idx] = result[position];
        result[position] = list[idx];
    }
    return result;
  });
  const shuffle = shuffler(Math.random);

  const mappedTickerItems =
    logos.map(item => (
      <TickerImage key={item} src={item} alt="logo" />
    ))

  return (
    <div>
      <Carousel />
      <Catalog isMainPage />
  
      <ClientsBlockWrapper>
        <ClientsBlockTitleWrapper>
          <BlockInfoComponent title="Клієнти" subtitle="Серед них" />
  
          <TickerWrapper onMouseOver={handleMouseOver} onMouseLeave={handleMouseLeave}>
            <Ticker speed={8} mode='await' move={isMouseOverTicker}>
              {() => (
                <ClientTicker>
                  {shuffle(mappedTickerItems)}
                </ClientTicker>
              )}
            </Ticker>
          </TickerWrapper>

          {/* <LogoContainer>{mappedTickerItems}</LogoContainer> */}
        </ClientsBlockTitleWrapper>
      </ClientsBlockWrapper>

      <AdvantagesBlock advantageItems={mockedAdvantages} />

      <AboutUsBlock />

      <ContactsAndFormBlock />
    </div>
  );
};

export default HomePageView;
