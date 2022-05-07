// node modules
import React from 'react';
// modules
import { ROUTES } from 'routing/registration';
// import { Link as RouterLink } from "react-router-dom";
import { colors } from 'utils/color';
// components
import {
  createTheme,
  ThemeProvider,
  CssBaseline,
} from '@mui/material';
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
    href: ROUTES.HOME,
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

const Header = () => (
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
              key: label,
              color: 'success',
              href,
              variant,
              size: 'small',
            }}
          >
            <StyledButtonText color={color}>
              {label}
            </StyledButtonText>
          </StyledButton>
        ))}
      </StyledStack>
    </StyledAppBar>
  </ThemeProvider>
);

export default Header;
