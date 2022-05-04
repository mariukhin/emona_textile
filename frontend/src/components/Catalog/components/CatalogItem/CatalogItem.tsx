// node modules
import React from 'react';
// components
import { Typography } from '@mui/material';
// styles
import {
  StyledPaper,
  StyledGrid,
  CatalogItemImageWrapper,
  CatalogButton,
} from '../../styles';

type CatalogItem = {
  title: CatalogData['title'],
  imageUrl: CatalogData['imageUrl'],
};

const CatalogItem: React.FC<CatalogItem> = ({ title, imageUrl }) => (
  <StyledGrid item xs={2} sm={4} md={4}>
    <StyledPaper>
      <CatalogItemImageWrapper
        theme={{
          main: imageUrl,
        }}
      >
        <CatalogButton color="success" variant="contained">
          <Typography
            textTransform="none"
            fontSize={20}
            fontWeight={700}
            sx={{ fontFamily: 'Comfortaa' }}
          >
            {title}
          </Typography>
        </CatalogButton>
      </CatalogItemImageWrapper>
    </StyledPaper>
  </StyledGrid>
);

export default CatalogItem;