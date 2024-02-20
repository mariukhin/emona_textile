import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Button, Toolbar, Stack, IconButton, Typography, List, ListItemButton, Link } from '@mui/material';
import {
  ContactBlockTextLink,
} from 'components/Footer/components/FooterInfoBlock/styles';

export const StyledToolbar = styled(Toolbar)`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  background-color: ${colors.background.white};
  box-shadow: none;

  @media ${device.mobile} {
    padding: 8px 16px;
  }

  @media ${device.tablet} {
    padding: 12px 24px;
  }

  @media ${device.laptopL} {
    padding: 20px 60px;
  }
`;

export const StyledLogo = styled(Link)`
  padding: 0;

  @media ${device.mobile} {
    img {
      width: 126px;
      height: 24px;
      margin-top: 8px;
    }
  }

  @media ${device.tablet} {
    img {
      width: 208px;
      height: 40px;
      margin-top: 0;
    }
  }
`;

export const StyledStack = styled(Stack)`
  display: flex;
  align-items: center;

  @media ${device.mobile} {
    flex-direction: row-reverse;

    a:nth-child(-n+4) {
      display: none;
    }

    a:last-child {
      margin-right: 30px;
    }
  }

  @media ${device.tablet} {
    a:last-child {
      margin-right: 35px;
    }
  }

  @media ${device.laptopL} {
    flex-direction: row;

    a:nth-child(-n+4) {
      display: block;
    }

    a:last-child {
      margin-right: 0;
    }
  }
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

  @media ${device.mobile} {
    height: 32px;
  }

  @media ${device.tablet} {
    height: 38px;
  }
`;

export const StyledButtonText = styled(Typography)`
  text-transform: uppercase;
  font-family: 'Nunito';
  font-weight: 700;

  @media ${device.mobile} {
    font-size: 14px;
    line-height: 24px;
  }

  @media ${device.tablet} {
    font-size: 16px;
    line-height: 26px;
  }
`;

export const StyledButtonTextDrawer = styled(Typography)`
  font-family: 'Comfortaa';
  font-weight: 700;

  @media ${device.mobile} {
    font-size: 24px;
    line-height: 36px;
  }
`;

export const StyledButtonTextAnd = styled(Typography)`
  font-family: 'Comfortaa';
  font-weight: 700;
  margin: 0 auto;
  padding-top: 30px;
  padding-bottom: 20px;

  @media ${device.mobile} {
    font-size: 18px;
    line-height: 28px;
  }
`;

export const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
  justifyContent: 'flex-start',
}));

export const StyledDrawerList = styled(List)`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
  width: 100%;
  padding-top: 50%;
`;

export const StyledDrawerButton = styled(ListItemButton)<{ href: string; }>`
  justify-content: center;
`;

export const StyledDrawerContactButton = styled(StyledButton)`
  max-width: 60%;
  margin: 0 auto;
  margin-top: 8px;

  @media ${device.mobile} {
    height: 48px;
  }
`;

export const StyledContactItemPhoneBlock = styled.div`
  display: flex;
  flex-direction: column;
  max-width: 60%;
  margin: 0 auto;
`;

export const StyledContactBlockTextLink = styled(ContactBlockTextLink)`
  padding-bottom: 12px;
`;

export const StyledButtonWrapper = styled.div`
  text-align: center;
  margin-top: 52px;
`;

export const StyledBurger = styled(IconButton)`
  color: ${colors.background.grey};
  display: block;
  margin-right: 0;

  @media ${device.mobile} {
    padding: 0 !important;
  }

  @media ${device.laptopL} {
    display: none;
  }
`;
