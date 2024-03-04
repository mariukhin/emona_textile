import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Grid } from '@mui/material';

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
  @media ${device.mobile} {
    height: 300px;
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
    height: 350px;
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

export const CatalogItemWrapper = styled.a`
  height: 100%;
  width: 100%;
  border-radius: 20px;
  box-shadow: 0px 4px 4px 0px #00000040;
  text-decoration: none;

  display: flex;
  justify-content: flex-start;
  flex-direction: column;
  transition: filter 300ms 100ms;

  :hover {
    cursor: pointer;
    filter: drop-shadow(0px 6px 32px rgba(0, 0, 0, 0.36));
  }

  :hover > div {
    background-color: ${colors.button.hover};
  }
`

export const CatalogItemImage = styled.img`
  width: 100%;
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;

  @media ${device.mobile} {
    min-height: 60%;
  }

  @media ${device.laptop} {
    min-height: 75%;
  }
`;

export const CatalogButton = styled.div`
  display: flex;
  flex: auto;
  justify-content: center;
  align-items: center;
  padding: 3px 0;
  background-color: ${colors.button.default};
  font-family: 'Comfortaa';
  font-size: 20px;
  line-height: 36px;
  text-align: center;
  font-style: normal;
  color: ${colors.text.white};
  border-bottom-left-radius: 20px;
  border-bottom-right-radius: 20px;
  width: 100%;
  transition: background-color 300ms 100ms;
  text-decoration: none;
`;
