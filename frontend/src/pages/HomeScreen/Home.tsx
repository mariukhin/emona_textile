// node modules
import React, { useState, useEffect } from 'react';
import Ticker from 'react-ticker';
// components
import AboutUsBlock from 'components/AboutUsBlock';
import AdvantagesBlock from 'components/AdvantagesBlock';
import BlockInfoComponent from 'components/BlockInfoComponent';
import { Carousel } from 'components/Carousel';
import Catalog from 'components/Catalog';
import { ContactsAndFormBlock } from 'components/ContactsAndFormBlock';
import { useStore } from 'modules/Stores';
// mocks
import { mockedCatalogItems } from 'components/Catalog/mocks';
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

  const { getCarouselItems } = useStore('CarouselStore');

  useEffect(() => {
    getCarouselItems();
  }, []);

  const handleMouseOver = () => {
    setIsMouseOverTicker(false);
  }

  const handleMouseLeave = () => {
    setIsMouseOverTicker(true);
  }

  return (
    <>
      <Carousel />
      <Catalog catalogItems={mockedCatalogItems} isMainPage />
  
      <ClientsBlockWrapper>
        <ClientsBlockTitleWrapper>
          <BlockInfoComponent title="Клієнти" subtitle="Серед них" />
  
          <TickerWrapper onMouseOver={handleMouseOver} onMouseLeave={handleMouseLeave}>
            <Ticker speed={8} mode='await' move={isMouseOverTicker}>
              {() => (
                <ClientTicker>
                  {logos.map(item => (
                    <TickerImage key={item} src={item} alt="logo" />
                  ))}
                </ClientTicker>
              )}
            </Ticker>
          </TickerWrapper>
        </ClientsBlockTitleWrapper>
      </ClientsBlockWrapper>

      <AdvantagesBlock advantageItems={mockedAdvantages} />

      <AboutUsBlock />

      <ContactsAndFormBlock />
    </>
  );
};

export default HomePageView;
