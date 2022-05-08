import styled from 'styled-components';
import { colors } from 'utils/color';
import { Button, AppBar, Stack, IconButton, Typography } from '@mui/material';

export const StyledAppBar = styled(AppBar)`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 22px 60px;
  background-color: ${colors.background.white};
  box-shadow: none;
`;

export const StyledLogo = styled(IconButton)`
  padding: 0;
`;

export const StyledStack = styled(Stack)`
  display: flex;
  flex-direction: row;
  align-items: center;
`;

export const StyledButton = styled(Button)`
  padding: 6px 18px;
  height: 38px;
  background-color: ${props => props.variant === 'text' ? 'inherit' : colors.button.default};

  span {
    margin-left: 4px;

    svg {
      color: ${colors.text.default};
      font-size: 24px !important;
    }
  }

  :hover{
    background-color: ${props => props.variant === 'text' ? 'inherit' : colors.button.default};
    p {
      color: ${props => props.variant === 'text' ? colors.text.orange : 'inherit'}
    }
  }
`;

export const StyledButtonText = styled(Typography)`
  font-size: 16px;
  line-height: 26px;
  text-transform: uppercase;
  font-family: 'Nunito';
  font-weight: 700;
`;

export const StyledButtonWrapper = styled.div`
  text-align: center;
  margin-top: 52px;
`;
