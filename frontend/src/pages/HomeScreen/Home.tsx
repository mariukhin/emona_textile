// node modules
import React, { useState } from 'react';
import Ticker from 'react-ticker';
// components
import AboutUsBlock from 'components/AboutUsBlock';
import AdvantagesBlock from 'components/AdvantagesBlock';
import BlockInfoComponent from 'components/BlockInfoComponent';
import { Carousel } from 'components/Carousel';
import { Catalog } from 'components/Catalog';
import { ContactsAndFormBlock } from 'components/ContactsAndFormBlock';
// import { useStore } from 'modules/Stores';
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
  LogoContainer,
} from './styles';

const HomePageView = () => {
  const [isMouseOverTicker, setIsMouseOverTicker] = useState(true);

  // const { getCarouselItems } = useStore('CarouselStore');
  // const { getCatalogItems } = useStore('CatalogStore');

  // useEffect(() => {
  //   getCarouselItems();
  //   getCatalogItems();
  // }, []);

  const handleMouseOver = () => {
    setIsMouseOverTicker(false);
  }

  const handleMouseLeave = () => {
    setIsMouseOverTicker(true);
  }

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
                  {mappedTickerItems}
                </ClientTicker>
              )}
            </Ticker>
          </TickerWrapper>

          <LogoContainer>{mappedTickerItems}</LogoContainer>
        </ClientsBlockTitleWrapper>
      </ClientsBlockWrapper>

      <AdvantagesBlock advantageItems={mockedAdvantages} />

      <AboutUsBlock />

      <ContactsAndFormBlock />
    </div>
  );
};

export default HomePageView;
