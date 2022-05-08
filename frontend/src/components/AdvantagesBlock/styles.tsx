import styled from 'styled-components';
import { colors } from 'utils/color';
import { Paper, Button, Typography } from '@mui/material';

export const AdvantagesBlockWrapper = styled.div`
  margin: 120px 0;
`;

export const StyledGridContainer = styled.div`
  margin: 0 auto;
  margin-top: 40px;
  display: flex;
  justify-content: space-between;
  width: 90%;
`;

export const StyledPaper = styled(Paper)`
  width: 24%;
  height: 254px;
  border-radius: 20px;
`;

export const BlockWrapper = styled.div`
  width: 100%;
  padding: 20px;
`;

export const BlockImage = styled.img`
  width: 80px;
  height: 80px;
  margin-bottom: 25px;
`;

export const BlockHeading = styled(Typography)`
  font-size: 24px;
  line-height: 30px;
  color: ${colors.text.default};
  margin-bottom: 6px;
  font-family: 'Comfortaa';
  font-weight: 700;
`;

export const BlockSubHeading = styled(Typography)`
  font-size: 16px;
  line-height: 20px;
  color: ${colors.text.greyLight};
  font-family: 'Montserrat';
  font-weight: 400;
`;

export const StyledButton = styled(Button)`
  padding: 10px 28px;
  background-color: ${colors.button.default};
`;
