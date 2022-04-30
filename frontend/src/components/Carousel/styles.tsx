import styled from 'styled-components';
import { colors } from 'utils/color';
import { Button, Typography, Stack, Fab } from '@mui/material';

export const CarouselContainer = styled.div`
  width: 100%;
  height: 100%;
  background-image: url(${props => props.theme.main});
`;

export const ContentWrapper = styled.div`
  width: 100%;
  height: 604px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

export const InfoBlock = styled.div`
  width: 80%;
  margin: 0 auto;
  height: 148px;
  padding-top: 220px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  text-align: center;
`;

export const CarouselHeading = styled(Typography)`
  font-size: 60px;
  line-height: 90px;
  color: ${colors.text.white};
  text-shadow: 0px 2px 24px #000000;
  margin-bottom: 20px;
`;

export const StyledButton = styled(Button)`
  padding: 6px 18px;
  height: 38px;
  background-color: ${colors.button.carousel};
  border-radius: 12px;
`;

export const StyledFab = styled(Fab)`
  background-color: #FFFFFF;
  opacity: 0.75;
`;

export const CarouselButtonsBlock = styled.div`
  display: flex;
  align-items: center;
  margin-top: 320px;
  justify-content: space-between;
`;

export const ItemsBlock = styled(Stack)`
  display: flex;
  margin: 0 20px;
`;

export const Item = styled.span`
  width: 28px;
  height: 6px;
  border-radius: 4px;
  background-color: ${(props) => props.theme.main};
`;
