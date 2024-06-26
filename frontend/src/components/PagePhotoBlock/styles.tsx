import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Button, Typography } from '@mui/material';

export const PagePhotoBlockContainer = styled.div<{ imageUrl: string; }>`
  height: 450px;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: top;
  background-image: url(${props => props.imageUrl});
  display: flex;
  justify-content: center;
  align-items: center;

  @media ${device.mobile} {
    height: 350px;
  }

  @media ${device.tablet} {
    height: 400px;
  }

  @media ${device.laptop} {
    height: 450px;
  }

  @media ${device.laptopL} {
    margin: 0 59px;
  }
`;

export const InfoBlock = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;

  @media ${device.mobile} {
    max-width: 85%;
  }

  @media ${device.tablet} {
    height: 128px;
    max-width: 85%;
  }

  @media ${device.laptop} {
    max-width: 70%;
  }

  @media ${device.laptopL} {
    height: 148px;
  }
`;

export const Heading = styled(Typography)`
  color: ${colors.text.white};
  text-shadow: 0px 2px 24px #000000;
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
  background-color: ${colors.button.carousel};
  border-radius: 12px;
  max-width: 85%;
`;

export const StyledButtonText = styled(Typography)`
  font-size: 20px;
  font-family: 'Montserrat';
  font-weight: 600;
  color: ${colors.text.white};
  text-transform: none;
`;
