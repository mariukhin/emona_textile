// node modules
import React from 'react';
import { colors } from 'utils/color';
import { observer } from 'mobx-react';
// modules
import CarouselService from '../service';
// components
import { Typography } from '@mui/material';
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

const Carousel = () => {
  const { carouselItems, currentItem } = useStore('CarouselStore');

  if (!carouselItems) return null;

  return (
    <CarouselContainer
      theme={{
        main: currentItem?.imageUrl,
      }}
    >
      <ContentWrapper>
        <InfoBlock>
          <CarouselHeading variant="body1" sx={{ fontFamily: 'Comfortaa' }}>
            {currentItem?.title}
          </CarouselHeading>
          <StyledButton color="warning" size="large" variant="contained">
            <Typography
              variant="button"
              textTransform="none"
              fontSize={20}
              fontWeight={600}
              sx={{ fontFamily: 'Montserrat' }}
            >
              {currentItem?.buttonText}
            </Typography>
          </StyledButton>
        </InfoBlock>
        <CarouselButtonsBlock>
          <StyledFab
            color="default"
            size="small"
            onClick={() => CarouselService.changeCurrentItem('Left')}
          >
            <ArrowBack sx={{ color: colors.background.green }} />
          </StyledFab>

          <ItemsBlock direction="row" spacing={1}>
            {carouselItems.map(({ id, isCurrent }) => (
              <Item
                key={id}
                theme={{
                  main: isCurrent
                    ? colors.button.default
                    : colors.background.default,
                }}
              />
            ))}
          </ItemsBlock>

          <StyledFab
            color="default"
            size="small"
            onClick={() => CarouselService.changeCurrentItem('Right')}
          >
            <ArrowForward sx={{ color: colors.background.green }} />
          </StyledFab>
        </CarouselButtonsBlock>
      </ContentWrapper>
    </CarouselContainer>
  );
};

export default observer(Carousel);
