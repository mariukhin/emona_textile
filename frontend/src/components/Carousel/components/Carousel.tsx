// node modules
import React from 'react';
import { colors } from 'utils/color';
import { observer } from 'mobx-react';
// modules
import CarouselService from '../service';
import { useStore } from 'modules/Stores';
// components
import { ArrowBack, ArrowForward } from '@mui/icons-material';
// styles
import {
  CarouselContainer,
  CarouselButtonsBlock,
  CarouselHeading,
  ContentWrapper,
  InfoBlock,
  StyledButton,
  StyledButtonText,
  ItemsBlock,
  StyledFab,
  Item,
} from '../styles';

const Carousel = () => {
  const { carouselItems, currentItem } = useStore('CarouselStore');

  if (!carouselItems) return null;

  return (
    <CarouselContainer imageUrl={ currentItem?.imageUrl || '' } backgroundColor={ currentItem?.backgroundColor || '' }>
      <ContentWrapper>
        <InfoBlock>
          <CarouselHeading>{currentItem?.title}</CarouselHeading>
          <StyledButton color="warning" size="large" variant="contained">
            <StyledButtonText>{currentItem?.btnText}</StyledButtonText>
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
