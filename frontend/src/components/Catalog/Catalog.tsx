// node modules
import React from 'react';
import { observer } from 'mobx-react';
// modules
import { useStore } from 'modules/Stores';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import CatalogItem from './components/CatalogItem';
// styles
import {
  CatalogWrapper,
  StyledGridContainer,
} from './styles';

type CatalogProps = {
  isMainPage?: boolean,
}

const Catalog: React.FC<CatalogProps> = ({ isMainPage = false }) => {
  const { catalogItems } = useStore('CatalogStore');

  return (
    <CatalogWrapper style={{ marginTop: isMainPage ? '120px' : '40px' }}>
      {isMainPage && <BlockInfoComponent title="Каталог" subtitle="Якісні товари" />}
  
      <StyledGridContainer container spacing={3}>
        {catalogItems && catalogItems.map((item) => (
          <CatalogItem key={item.id} title={item.title} imageUrl={item.imageUrl} />
        ))}
      </StyledGridContainer>
    </CatalogWrapper>
  )
};

export default observer(Catalog);
