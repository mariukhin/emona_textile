// node modules
import React, { useId } from 'react';
import { observer } from 'mobx-react';
// modules
import { ROUTES } from 'routing/registration';
import { colors } from 'utils/color';
import { useStore } from 'modules/Stores';
// components
// import { KeyboardArrowDown } from '@mui/icons-material';
import MenuIcon from '@mui/icons-material/Menu';
// styles
import {
  StyledAppBar,
  StyledLogo,
  StyledStack,
  StyledButton,
  StyledButtonText,
  StyledBurger,
} from './styles';

const headersData: HeadersData[] = [
  {
    label: 'Головна',
    href: ROUTES.HOME,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: 'Каталог',
    href: ROUTES.CATALOG,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: 'Про нас',
    href: ROUTES.ABOUT,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: 'Зв’язатися',
    href: ROUTES.HOME,
    variant: 'contained',
    color: colors.text.white,
  },
];

const Header = () => {
  const { isOnRoute } = useStore('RoutingStore');

  const getItemColor = (href: string, itemColor: string) => {
    if (href !== ROUTES.HOME && isOnRoute(href)) {
      return colors.text.orange;
    } else {
      return itemColor;
    }
  }

  return (
    <div>  
      <StyledAppBar position="fixed">
        <StyledLogo>
          <img src='assets/logo.svg' alt="Emona logo" />
        </StyledLogo>
  
        <StyledStack direction="row" spacing={1}>
          <StyledBurger
            size="large"
            edge="start"
            aria-label="menu"
            sx={{
              mt: 1,
              mr: 2,
            }}
          >
            <MenuIcon fontSize="large" />
          </StyledBurger>

          {headersData.map(({ label, href, variant, color }) => (
            <StyledButton
              {...{
                key: useId(),
                color: 'success',
                href,
                variant,
                size: 'small',
                // endIcon:  
                //   label === 'Каталог' && 
                //     <KeyboardArrowDown
                //       fontSize="large"
                //     />,
              }}
            >
              <StyledButtonText color={getItemColor(href, color)}>
                {label}
              </StyledButtonText>
            </StyledButton>
          ))}
        </StyledStack>
      </StyledAppBar>
    </div>
  );
};

export default observer(Header);
