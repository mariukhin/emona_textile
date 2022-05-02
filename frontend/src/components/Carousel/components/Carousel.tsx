// node modules
import * as R from 'ramda';
import React from 'react';
import { colors } from 'utils/color';
import { observer } from 'mobx-react';
// modules
import CarouselService from '../service';
// components
import { Typography, createTheme, ThemeProvider } from '@mui/material';
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
} from '../styles';
import { useStore } from 'modules/Stores';

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

const Carousel = () => {
  const { carouselItems } = useStore('CarouselStore');

  if (!carouselItems) return null;

  const currentItem = R.find((item) => item.isCurrent, carouselItems);

  return (
    <CarouselContainer
      theme={{
        main: currentItem?.imageUrl,
      }}
    >
      <ContentWrapper>
        <InfoBlock>
          <ThemeProvider theme={HeadingTheme}>
            <CarouselHeading variant="h1">{currentItem?.title}</CarouselHeading>
          </ThemeProvider>
          <ThemeProvider theme={ButtonTheme}>
            <StyledButton color="warning" size="large" variant="contained">
              <Typography
                variant="button"
                textTransform="none"
                fontSize={20}
                fontWeight={600}
              >
                {currentItem?.buttonText}
              </Typography>
            </StyledButton>
          </ThemeProvider>
        </InfoBlock>
        <CarouselButtonsBlock>
          <StyledFab color="default" size="small" onClick={() => CarouselService.changeCurrentItem('Left')}>
            <ArrowBack sx={{ color: colors.background.green }} />
          </StyledFab>

          <ItemsBlock direction="row" spacing={1}>
            {carouselItems.map(({ id, isCurrent }) => (
              <Item
                key={id}
                theme={{
                  main: isCurrent
                    ? colors.button.default
                    : colors.background.grey,
                }}
              />
            ))}
          </ItemsBlock>

          <StyledFab color="default" size="small" onClick={() => CarouselService.changeCurrentItem('Right')}>
            <ArrowForward sx={{ color: colors.background.green }} />
          </StyledFab>
        </CarouselButtonsBlock>
      </ContentWrapper>
    </CarouselContainer>
  );
};

export default observer(Carousel);
