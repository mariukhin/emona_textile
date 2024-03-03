// node modules
import React from 'react';
// styles
import {
  StyledGrid,
  CatalogItemWrapper,
  CatalogItemImage,
  CatalogButton,
} from '../../styles';
import { ROUTES } from 'routing/registration';

type CatalogItem = {
  title: CatalogData['title'],
  imageUrl: CatalogData['imageUrl'],
};

const CatalogItem: React.FC<CatalogItem> = ({ title, imageUrl }) => (
  <StyledGrid item xs={12} sm={6} md={6} lg={4}>
    <CatalogItemWrapper href={`${ROUTES.CATALOG_ITEM}?title=${title}`}>
      <CatalogItemImage src={ imageUrl } alt='catalog item image'/>
      <CatalogButton>{title}</CatalogButton>
    </CatalogItemWrapper>
  </StyledGrid>
); 

export default CatalogItem;
