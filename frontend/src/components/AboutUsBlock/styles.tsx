import styled from 'styled-components';
import { colors } from 'utils/color';
import { Paper, Button, Typography } from '@mui/material';

export const AboutUsBlockWrapper = styled.div`
  margin: 0 0 120px;
  position: relative;
`;

export const AboutUsPhotoBlock = styled.div`
  width: 832px;
  height: 734px;
  background-image: url('assets/about-us.png');
  border-radius: 0 20px 20px 0;
  background-size: cover;
  background-position: center;
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
  width: 669px;
  height: 636px;
  border-radius: 20px;
  box-shadow: 0px 6px 24px rgba(0, 0, 0, 0.12);
  background-color: #fff;
  left: 708px;
  top: 52px;
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
  width: 90%;
  justify-content: space-between;
  margin-bottom: 44px;
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
