import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Typography } from '@mui/material';

export const BlockHeading = styled(Typography)`
  line-height: 52px;
  color: ${colors.text.greyDark};
  margin-bottom: 16px;
  font-family: 'Comfortaa';
  font-weight: 700;

  @media ${device.mobile} {
    font-size: 36px;
  }

  @media ${device.laptop} {
    font-size: 48px;
  }
`;

export const BlockSubHeading = styled(Typography)`
  font-size: 20px;
  line-height: 22px;
  color: ${colors.text.green};
  font-family: 'Montserrat';
  font-weight: 700;
`;

