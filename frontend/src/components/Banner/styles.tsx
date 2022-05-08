import styled from 'styled-components';
import { colors } from 'utils/color';
import { Typography } from '@mui/material';

export const BannerContainer = styled.div`
  border: 2px solid ${colors.background.green};
  border-radius: 20px;
  height: 106px;
`;

export const BannerWrapper = styled.div`
  padding: 16px 20px;
  display: flex;
  height: 100%;
  flex-direction: column;
  justify-content: space-between;
`;

export const BannerHeading = styled(Typography)`
  font-size: 36px;
  line-height: 40px;
  color: ${colors.text.greyDark};
  font-family: 'Comfortaa';
  font-weight: 700;
`;

export const BannerSubHeading = styled(Typography)`
  font-size: 16px;
  line-height: 20px;
  color: ${colors.text.greyLight};
  font-family: 'Montserrat';
  font-weight: 400;
`;

