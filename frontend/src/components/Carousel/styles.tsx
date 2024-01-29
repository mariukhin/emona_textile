import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Button, Typography, Stack, Fab } from '@mui/material';

export const CarouselContainer = styled.div<{ imageUrl: string; backgroundColor: string; }>`
  height: 100%;
  background-image: url(${props => props.imageUrl});
  background-repeat: no-repeat;
  background-size: cover;
  background-position: top;

  @media ${device.mobile} {
    margin: 0 12px;
  }

  @media ${device.tablet} {
    margin: 0 22px;
  }

  @media ${device.laptopL} {
    margin: 0 59px;
  }
`;

export const ContentWrapper = styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  @media ${device.mobile} {
    height: 496px;
  }

  @media ${device.laptopL} {
    height: 604px;
  }
`;

export const InfoBlock = styled.div`
  margin: 0 auto;
  height: 148px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  text-align: center;

  @media ${device.mobile} {
    width: 90%;
    padding-top: 100px;
  }

  @media ${device.laptop} {
    width: 95%;
    padding-top: 150px;
  }

  @media ${device.laptopL} {
    width: 90%;
    padding-top: 220px;
  }
`;

export const CarouselHeading = styled(Typography)`
  color: ${colors.text.white};
  text-shadow: 0px 2px 24px #000000;
  margin-bottom: 20px;
  font-family: 'Comfortaa';

  @media ${device.mobile} {
    font-size: 36px;
    line-height: 54px;
  }

  @media ${device.tablet} {
    font-size: 48px;
    line-height: 72px;
  }

  @media ${device.laptopL} {
    font-size: 60px;
    line-height: 90px;
  }
`;

export const StyledButton = styled(Button)`
  padding: 6px 18px;
  height: 38px;
  background-color: ${colors.button.carousel};
  border-radius: 24px;
`;

export const StyledButtonText = styled(Typography)`
  font-size: 20px;
  font-family: 'Montserrat';
  font-weight: 600;
  color: ${colors.text.white};
  text-transform: none;
`;

export const StyledFab = styled(Fab)`
  background-color: #FFFFFF;
  opacity: 0.75;
`;

export const CarouselButtonsBlock = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;

  @media ${device.mobile} {
    margin-top: 270px;
  }

  @media ${device.tablet} {
    margin-top: 250px;
  }

  @media ${device.laptopL} {
    margin-top: 320px;
  }
`;

export const ItemsBlock = styled(Stack)`
  display: flex;
  margin: 0 20px;
`;

export const Item = styled.span`
  width: 28px;
  height: 6px;
  border-radius: 4px;
  background-color: ${(props) => props.theme.main};
`;
