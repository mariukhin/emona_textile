import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Grid, Paper, Button, Typography } from '@mui/material';

export const CatalogWrapper = styled.div<{ isMainPage: boolean; }>`
  @media ${device.mobile} {
    margin: 60px 12px;
    margin-top: ${props => !props.isMainPage && '40px'} !important;
  }

  @media ${device.tablet} {
    margin: 80px 22px;
  }

  @media ${device.laptopL} {
    margin: 120px 59px;
  }
`;

export const StyledGridContainer = styled(Grid)`
  margin: 0 auto;
  margin-top: 20px;
  width: 100%;
`;

export const StyledGrid = styled(Grid)`
  padding-top: 0 !important;
  @media ${device.mobile} {
    height: 272px;
    padding-left: 0 !important;
  }

  @media ${device.tablet} {
    padding-left: 18px !important;

    :nth-of-type(2n+1) {
      padding-left: 0 !important;
    }
  }

  @media ${device.laptop} {
    padding-left: 24px !important;
    height: 248px;
  }

  @media ${device.laptopL} {
    :nth-of-type(3n+1) {
      padding-left: 0 !important;
    }

    :nth-of-type(2n+1):not(:first-of-type) {
      padding-left: 24px !important;
    }
  }
`;

export const StyledPaper = styled(Paper)`
  height: 100%;
  width: 100%;
  border-radius: 20px;
  box-shadow: 0px 4px 4px 0px #00000040;
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
  background-position: center;
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
  padding-left: 20px;
`;

export const ItemInfoDescriptionListItem = styled.li`
  font-family: 'Montserrat';
  font-style: regular;
  font-weight: 400;
  font-size: 16px;
  line-height: 24px;
  letter-spacing: 3%;
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
