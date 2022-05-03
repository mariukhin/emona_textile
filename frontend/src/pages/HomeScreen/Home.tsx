// node modules
import React from 'react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import { Carousel } from 'components/Carousel';
import Catalog from 'components/Catalog';
// mocks
import { mockedCatalogItems } from 'components/Catalog/mocks';
// styles
import { ClientsBlockWrapper, ClientsBlockTitleWrapper } from './styles';

const HomePageView = () => (
  <>
    <Carousel />
    <Catalog catalogItems={mockedCatalogItems} />

    <ClientsBlockWrapper>
      <ClientsBlockTitleWrapper>
        <BlockInfoComponent title="Клієнти" subtitle="Серед них" />
      </ClientsBlockTitleWrapper>
    </ClientsBlockWrapper>
  </>
);

export default HomePageView;
