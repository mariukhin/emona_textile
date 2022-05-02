// node modules
import React from 'react';
// components
import { Carousel } from 'components/Carousel';
import BlockInfoComponent from 'components/BlockInfoComponent';
import { Grid, createTheme } from '@mui/material';
// styles
import {
  CatalogueWrapper,
  StyledPaper,
  StyledStack
} from './styles';

const HeadingTheme = createTheme({
  typography: {
    fontFamily: 'Comfortaa',
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: `
          @font-face {
            font-family: 'Comfortaa';
            src: url('fonts/Comfortaa/Comfortaa-Bold.ttf');
            font-weight: 700;
          }
        `,
    },
  },
});

const catalogueItems = [
  { id: 1232 },
  { id: 1211 },
  { id: 1542 },
  { id: 3488 },
  { id: 2499 },
  { id: 4539 },
];

const HomePageView = () => (
  <>
    <Carousel />
    <CatalogueWrapper>
      <BlockInfoComponent title='Каталог' subtitle='Якісні товари'/>

      <StyledStack container spacing={{ xs: 3, md: 3 }} columns={{ xs: 4, sm: 8, md: 12 }}>
        {catalogueItems.map(item => (
          <Grid item xs={2} sm={4} md={4} key={item.id}>
            <StyledPaper elevation={3} />
          </Grid>
        ))}
      </StyledStack>
    </CatalogueWrapper>
  </>
);

export default HomePageView;
