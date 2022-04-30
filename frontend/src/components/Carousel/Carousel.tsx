// node modules
import React from 'react';
import { colors } from 'utils/color';
// components
import {
  Typography,
  createTheme,
  ThemeProvider,
} from '@mui/material';
import { ArrowBack, ArrowForward } from '@mui/icons-material';
import {
  CarouselContainer,
  CarouselButtonsBlock,
  CarouselHeading,
  ContentWrapper,
  InfoBlock,
  StyledButton,
  ItemsBlock,
  StyledFab,
  Item,
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

const ButtonTheme = createTheme({
  typography: {
    fontFamily: 'Montserrat',
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: `
        @font-face {
          font-family: 'Montserrat';
          src: url('fonts/Montserrat/Montserrat-SemiBold.ttf');
          font-weight: 600;
        }
      `,
    },
  },
});

const items = [
  { id: '1', isCurrent: false },
  { id: '2', isCurrent: false },
  { id: '3', isCurrent: true },
  { id: '4', isCurrent: false },
]


const Carousel = () => (
  <CarouselContainer>
    <ContentWrapper>
      <InfoBlock>
        <ThemeProvider theme={HeadingTheme}>
          <CarouselHeading variant="h1">
            Готельно-рестораний текстиль
          </CarouselHeading>
        </ThemeProvider>
        <ThemeProvider theme={ButtonTheme}>
          <StyledButton
            color="warning"
            size="large"
            variant="contained"
          >
            <Typography variant='button' textTransform="none" fontSize={20} fontWeight={600}>
              Сауна, СПА, басейн
            </Typography>
          </StyledButton>
        </ThemeProvider>
      </InfoBlock>
      <CarouselButtonsBlock>
        <StyledFab color="default" size="small">
          <ArrowBack sx={{ color: colors.background.green }} />
        </StyledFab>

        <ItemsBlock direction="row" spacing={1}>
          {items.map(({ id, isCurrent }) =>
            <Item
              key={id}
              theme={{ 
                main: isCurrent ? colors.button.default : colors.background.grey
              }} 
            /> 
          )}
        </ItemsBlock>

        <StyledFab color="default" size="small">
          <ArrowForward sx={{ color: colors.background.green }} />
        </StyledFab>
      </CarouselButtonsBlock>
    </ContentWrapper>
  </CarouselContainer>
);

export default Carousel;
