// node modules
import React, { useId } from 'react';
import { observer } from 'mobx-react';
// modules
import { ROUTES } from 'routing/registration';
import { colors } from 'utils/color';
import { useStore } from 'modules/Stores';
// components
import {
  createTheme,
  ThemeProvider,
  CssBaseline,
} from '@mui/material';
import { KeyboardArrowDown } from '@mui/icons-material';
// styles
import {
  StyledAppBar,
  StyledLogo,
  StyledStack,
  StyledButton,
  StyledButtonText,
} from './styles';

const theme = createTheme({
  components: {
    MuiCssBaseline: {
      styleOverrides: `
        @font-face {
          font-family: 'Nunito';
          src: url('fonts/Nunito/Nunito-Bold.ttf');
          font-weight: 700;
        }

        @font-face {
          font-family: 'Nunito';
          src: url('fonts/Nunito/Nunito-Regular.ttf');
          font-weight: 400;
        }

        @font-face {
          font-family: 'Comfortaa';
          src: url('fonts/Comfortaa/Comfortaa-Bold.ttf');
          font-weight: 700;
        }

        @font-face {
          font-family: 'Montserrat';
          src: url('fonts/Montserrat/Montserrat-SemiBold.ttf');
          font-weight: 600;
        }

        @font-face {
          font-family: 'Montserrat';
          src: url('fonts/Montserrat/Montserrat-Regular.ttf');
          font-weight: 400;
        }
      `,
    },
  },
});

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
    href: ROUTES.HOME,
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
    <ThemeProvider theme={theme}>
      <CssBaseline />
  
      <StyledAppBar position="fixed">
        <StyledLogo>
          <img src='assets/logo.svg' alt="Emona logo" />
        </StyledLogo>
  
        <StyledStack direction="row" spacing={1}>
          {headersData.map(({ label, href, variant, color }) => (
            <StyledButton
              {...{
                key: useId(),
                color: 'success',
                href,
                variant,
                size: 'small',
                endIcon:  
                  label === 'Каталог' && 
                    <KeyboardArrowDown
                      fontSize="large"
                    />,
              }}
            >
              <StyledButtonText color={getItemColor(href, color)}>
                {label}
              </StyledButtonText>
            </StyledButton>
          ))}
        </StyledStack>
      </StyledAppBar>
    </ThemeProvider>
  );
};

export default observer(Header);
