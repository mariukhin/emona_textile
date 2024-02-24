import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Typography } from '@mui/material';

export const AboutUsWrapper = styled.div`
  margin: 0 auto;

  @media ${device.mobile} {
    width: 95%;
    padding: 80px 0 120px;
  }

  @media ${device.laptopL} {
    width: 60%;
    padding: 80px 0 140px;
  }
`;

export const AboutUsText = styled(Typography)`
  font-size: 16px;
  line-height: 24px;
  font-weight: 400;
  color: ${colors.text.greyLight};
  font-family: 'Montserrat';
`;

export const BannerContainer = styled.div`
  display: flex;
  flex-wrap: wrap;

  @media ${device.mobile} {
    height: 220px;
    justify-content: space-evenly;
    align-content: space-between;
    width: 95%;
    margin: 0 auto;
    margin-top: 60px;
    margin-bottom: 60px;
  }

  @media ${device.tablet} {
    justify-content: center;
    gap: 20px;
    height: 106px;
    width: 100%;
  }
`;

export const BlockImage = styled.img`
  width: 100%;
  height: 420px;
`;