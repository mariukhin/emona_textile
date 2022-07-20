import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Grid, Paper, Button } from '@mui/material';

export const CatalogWrapper = styled.div`
  @media ${device.mobile} {
    margin: 100px 12px;
  }

  @media ${device.tablet} {
    margin: 120px 22px;
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
    height: 312px;
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
  position: relative;
  height: 100%;
  width: 100%;
  border-radius: 20px;
  box-shadow: 0px 4px 4px 0px #00000040;
`;

export const CatalogItemImageWrapper = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  margin-right: -50%;
  transform: translate(-50%, -50%);
  width: 98%;
  height: 97%;
  border-radius: 20px;
  background-image: url(${props => props.theme.main});
  background-size: cover;
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const CatalogButton = styled(Button)`
  max-width: 75%;
  background-color: ${colors.button.default};
  border-radius: 12px;
  opacity: 0.8;

  span:first-of-type {
    display: none;
  }

  :hover span:first-of-type {
    margin-top: 8px;
    display: block;
    transition: display 5s;
  }
`;
