// node modules
import React from 'react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import CatalogItem from './components/CatalogItem';
// styles
import {
  CatalogWrapper,
  StyledGridContainer,
} from './styles';

type CatalogProps = {
  catalogItems: CatalogData[],
}

const Catalog: React.FC<CatalogProps> = ({ catalogItems }) => (
  <CatalogWrapper>
    <BlockInfoComponent title="Каталог" subtitle="Якісні товари" />

    <StyledGridContainer container spacing={3}>
      {catalogItems.map((item) => (
        <CatalogItem key={item.id} title={item.title} imageUrl={item.imageUrl} />
      ))}
    </StyledGridContainer>
  </CatalogWrapper>
);

export default Catalog;
