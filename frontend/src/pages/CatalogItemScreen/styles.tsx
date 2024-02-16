import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Grid, Paper, Button, Typography } from '@mui/material';

export const CatalogWrapper = styled.div<{ isMainPage: boolean; }>`
  @media ${device.mobile} {
    margin: 60px 12px;
  }

  @media ${device.tablet} {
    margin: 80px 22px;
  }

  @media ${device.laptopL} {
    margin: 46px 60px 120px 37px;
  }
`;

export const StyledGridContainer = styled(Grid)`
  margin: 0 auto;
  width: 100%;
`;

export const StyledGrid = styled(Grid)`
  max-height: 272px;
`;

export const StyledPaper = styled(Paper)`
  height: 100%;
  width: 100%;
  border-radius: 20px;
  box-shadow: 6px 0px 24px 0px rgba(0,0,0,0.12);
`;

export const ItemContainer = styled.div`
  padding: 24px;
  display: flex;
  justify-content: flex-start;
  gap: 24px;
`;

export const ItemImage = styled.div`
  width: 200px;
  height: 200px;
  border-radius: 20px;
  background-image: url(${props => props.theme.main});
  background-size: cover;
`;

export const ItemInfoBlock = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  max-height: 200px;
`;

export const ItemInfoBlockTitle = styled(Typography)`
  font-family: 'Comfortaa';
  font-weight: 700;
  font-size: 24px;
  color: ${colors.text.greyDark};
`;

export const ItemInfoDescriptionList = styled.ul`
  margin: 0;
  padding-left: 25px;
`;

export const ItemInfoDescriptionListItem = styled.li`
  font-family: 'Montserrat';
  font-style: regular;
  font-weight: 400;
  font-size: 16px;
  line-height: 24px;
  letter-spacing: 3%;
  color: ${colors.text.greyDark};
`;

export const ItemButton = styled(Button)`
  display: flex;
  max-width: 117px;

  span:first-of-type {
    display: block;
    line-height: 0.7;
  }
`;

export const ItemButtonText = styled(Typography)`
  font-family: 'Nunito';
  font-weight: 700;
  font-size: 13px;
  line-height: 22px;
  letter-spacing: 1px;
`;
