import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Button, Typography } from '@mui/material';

export const PagePhotoBlockContainer = styled.div`
  width: 100%;
  height: 366px;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: top;
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const InfoBlock = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  text-align: center;

  @media ${device.mobile} {
    height: 110px;
  }

  @media ${device.tablet} {
    height: 128px;
    width: 35%;
  }

  @media ${device.laptop} {
    width: 30%;
  }

  @media ${device.laptopL} {
    height: 148px;
  }
`;

export const Heading = styled(Typography)`
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
  border-radius: 12px;
`;

export const StyledButtonText = styled(Typography)`
  font-size: 20px;
  font-family: 'Montserrat';
  font-weight: 600;
  color: ${colors.text.white};
  text-transform: none;
`;
