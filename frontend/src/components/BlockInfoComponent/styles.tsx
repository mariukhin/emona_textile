import styled from 'styled-components';
import { colors } from 'utils/color';
import { Typography } from '@mui/material';

export const BlockContainer = styled.div`
  text-align: center;  
`

export const BlockHeading = styled(Typography)`
  font-size: 48px;
  line-height: 52px;
  color: ${colors.text.default};
  margin-bottom: 10px;
`;

export const BlockSubHeading = styled(Typography)`
  font-size: 20px;
  line-height: 22px;
  color: ${colors.text.green};
`;

