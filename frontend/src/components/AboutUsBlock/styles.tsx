import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Paper, Button, Typography } from '@mui/material';

export const AboutUsBlockWrapper = styled.div`
  position: relative;

  @media ${device.mobile} {
    margin: 0 0 1820px;
  }

  @media ${device.tablet} {
    margin: 0 0 800px;
  }

  @media ${device.laptop} {
    margin: 0 0 450px;
  }

  @media ${device.laptopL} {
    margin: 0 0 400px;
  }
`;

export const AboutUsPhotoBlock = styled.div`
  height: 734px;
  background-image: url('assets/about-us.png');
  border-radius: 0 20px 20px 0;
  background-size: cover;

  @media ${device.mobile} {
    width: 359px;
    height: 317px;
    background-position: right;
  }

  @media ${device.tablet} {
    width: 644px;
    height: 569px;
  }

  @media ${device.laptop} {
    width: 832px;
    height: 734px;
  }

  @media ${device.laptopL} {
    background-position: center;
  }
`;

export const StyledGridContainer = styled.div`
  margin: 0 auto;
  margin-top: 40px;
  display: flex;
  justify-content: space-between;
  width: 90%;
`;

export const StyledPaper = styled(Paper)`
  position: absolute;
  border-radius: 20px;
  box-shadow: 0px 6px 24px rgba(0, 0, 0, 0.12);
  background-color: #fff;

  @media ${device.mobile} {
    left: 16px;
    top: 159px;
    width: 359px;
    height: 1900px;
  }

  @media ${device.tablet} {
    left: 124px;
    top: 249px;
    width: 620px;
    height: 1060px;
  }

  @media ${device.laptop} {
    left: 397px;
    top: 52px;
    width: 610px;
    height: 1080px;
  }

  @media ${device.laptopL} {
    left: 620px;
    top: 52px;
    width: 757px;
    height: 940px;
  }
`;

export const PaperWrapper = styled.div`
  padding: 32px 40px;
`;

export const BlockText = styled(Typography)`
  font-size: 16px;
  line-height: 24px;
  color: ${colors.text.greyLight};
  margin: 28px 0;
  font-family: 'Montserrat';
  font-weight: 400;
`;

export const BannerContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;

  @media ${device.mobile} {
    width: 100%;
    margin-bottom: 20px;
  }

  @media ${device.tablet} {
    margin-bottom: 44px;
  }

  @media ${device.laptopL} {
    width: 80%;
  }
`;

export const StyledButton = styled(Button)`
  padding: 10px 20px 10px 28px;
  background-color: ${colors.button.carousel};
`;

export const StyledButtonText = styled(Typography)`
  font-size: 16px;
  line-height: 26px;
  text-transform: uppercase;
  font-family: 'Nunito';
  font-weight: 700;
`;
