import styled from 'styled-components';
import { colors } from 'utils/color';
import { Typography } from '@mui/material';

export const AboutUsWrapper = styled.div`
  width: 60%;
  margin: 0 auto;
  padding: 80px 0 140px;
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
  width: 100%;
  justify-content: space-between;
  margin: 60px 0;
`;

export const BlockImage = styled.img`
  width: 100%;
  height: 420px;
`;