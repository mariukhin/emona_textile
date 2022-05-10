// node modules
import React, { useState } from 'react';
import Ticker from 'react-ticker';
// components
import AboutUsBlock from 'components/AboutUsBlock';
import AdvantagesBlock from 'components/AdvantagesBlock';
import BlockInfoComponent from 'components/BlockInfoComponent';
import { Carousel } from 'components/Carousel';
import Catalog from 'components/Catalog';
import { ContactsAndFormBlock } from 'components/ContactsAndFormBlock';
// mocks
import { mockedCatalogItems } from 'components/Catalog/mocks';
import { mockedAdvantages } from 'components/AdvantagesBlock/mocks';
// styles
import {
  ClientsBlockWrapper,
  ClientsBlockTitleWrapper,
  ClientTicker,
  TickerWrapper,
  TickerImage,
} from './styles';

const logos = [
  'assets/ClientLogos/image-1.svg',
  'assets/ClientLogos/image-2.svg',
  'assets/ClientLogos/image-3.svg',
  'assets/ClientLogos/image-4.svg',
  'assets/ClientLogos/image-5.svg',
  'assets/ClientLogos/image-6.svg',
  'assets/ClientLogos/image-7.svg',
  'assets/ClientLogos/image-8.svg',
  'assets/ClientLogos/image-9.svg',
  'assets/ClientLogos/image-10.svg',
  'assets/ClientLogos/image-11.svg',
  'assets/ClientLogos/image-12.svg',
  'assets/ClientLogos/image-13.svg',
  'assets/ClientLogos/image-14.svg',
  'assets/ClientLogos/image-15.svg',
  'assets/ClientLogos/image-16.svg',
  'assets/ClientLogos/image-17.svg',
  'assets/ClientLogos/image-18.svg',
  'assets/ClientLogos/image-19.svg',
  'assets/ClientLogos/image-20.svg',
  'assets/ClientLogos/image-21.svg',
  'assets/ClientLogos/image-22.svg',
  'assets/ClientLogos/image-23.svg',
  'assets/ClientLogos/image-24.svg',
];

const HomePageView = () => {
  const [isMouseOverTicker, setIsMouseOverTicker] = useState(true);

  const handleMouseOver = () => {
    setIsMouseOverTicker(false);
  }

  const handleMouseLeave = () => {
    setIsMouseOverTicker(true);
  }

  return (
    <>
      <Carousel />
      <Catalog catalogItems={mockedCatalogItems} />
  
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
