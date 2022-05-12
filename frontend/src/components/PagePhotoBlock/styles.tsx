import styled from 'styled-components';
import { colors } from 'utils/color';
import { Button, Typography } from '@mui/material';

export const PagePhotoBlockContainer = styled.div`
  width: 100%;
  height: 366px;
  background-repeat: no-repeat;
  background-size: cover;
  display: flex;
  justify-content: flex-end;
  align-items: center;
`;

export const InfoBlock = styled.div`
  width: 20%;
  height: 148px;
  display: flex;
  margin-right: 15%;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  text-align: center;
`;

export const Heading = styled(Typography)`
  font-size: 60px;
  line-height: 90px;
  color: ${colors.text.white};
  text-shadow: 0px 2px 24px #000000;
  margin-bottom: 20px;
  font-family: 'Comfortaa';
`;

export const StyledButton = styled(Button)`
  padding: 6px 18px;
  height: 38px;
  background-color: ${colors.button.carousel};
  border-radius: 12px;
`;

export const StyledButtonText = styled(Typography)`
  font-size: 20px;
  font-family: 'Montserrat';
  font-weight: 600;
  color: ${colors.text.white};
  text-transform: none;
`;
