// node modules
import React from 'react';
// modules
// import { useStore } from 'modules/Stores';
// components
import { Typography } from '@mui/material';
import { ArrowForward } from '@mui/icons-material';
// styles
import {
  StyledPaper,
  StyledGrid,
  CatalogItemImageWrapper,
  CatalogButton,
} from '../../styles';
import { ROUTES } from 'routing/registration';

type CatalogItem = {
  title: CatalogData['title'],
  imageUrl: CatalogData['imageUrl'],
};

const CatalogItem: React.FC<CatalogItem> = ({ title, imageUrl }) => {
  // const { setCurrentCatalogItem } = useStore('CatalogItemStore');

  // const onCatalogButtonClick = () => {
  //   setCurrentCatalogItem(title);
  //   console.log('Here');
  // };

  return (
    <StyledGrid item xs={12} sm={6} md={6} lg={4}>
      <StyledPaper>
        <CatalogItemImageWrapper
          theme={{
            main: imageUrl,
          }}
        >
          <CatalogButton href={`${ROUTES.CATALOG_ITEM}?title=${title}`} color="success" variant="contained" endIcon={<ArrowForward />}>
            <Typography
              textTransform="none"
              fontSize={20}
              fontWeight={700}
              sx={{ fontFamily: 'Comfortaa', textAlign: 'center' }}
            >
              {title}
            </Typography>
          </CatalogButton>
        </CatalogItemImageWrapper>
      </StyledPaper>
    </StyledGrid>
  ); 
}

export default CatalogItem;
