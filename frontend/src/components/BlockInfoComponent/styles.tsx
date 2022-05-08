import styled from 'styled-components';
import { colors } from 'utils/color';
import { Typography } from '@mui/material';

export const BlockHeading = styled(Typography)`
  font-size: 48px;
  line-height: 52px;
  color: ${colors.text.default};
  margin-bottom: 16px;
  font-family: 'Comfortaa';
  font-weight: 700;
`;

export const BlockSubHeading = styled(Typography)`
  font-size: 20px;
  line-height: 22px;
  color: ${colors.text.green};
  font-family: 'Montserrat';
  font-weight: 700;
`;

