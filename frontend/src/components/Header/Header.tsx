// node modules
import React from 'react';
// modules
import { ROUTES } from 'routing/registration';
// import { Link as RouterLink } from "react-router-dom";
import styled from 'styled-components';
import { colors } from 'utils/color';
// components
import {
  AppBar,
  Button,
  Typography,
  IconButton,
  Stack,
  createTheme,
  ThemeProvider,
  CssBaseline,
} from '@mui/material';

const theme = createTheme({
  // typography: {
  //   fontFamily: 'Nunito',
  // },
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

const StyledAppBar = styled(AppBar)`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 20px 60px;
  background-color: ${colors.background.white};
  box-shadow: none;
`;

const StyledLogo = styled(IconButton)`
  padding: 0;
`;

const StyledStack = styled(Stack)`
  display: flex;
  flex-direction: row;
  align-items: center;
`;

const StyledButton = styled(Button)`
  padding: 6px 18px;
  height: 38px;
  background-color: ${props => props.variant === 'text' ? 'inherit' : colors.button.default};
`;

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
            <Typography
              variant="body1"
              fontSize="16px"
              lineHeight="1"
              color={color}
              textTransform="uppercase"
              sx={{
                fontFamily: 'Nunito',
              }}
            >
              {label}
            </Typography>
          </StyledButton>
        ))}
      </StyledStack>
    </StyledAppBar>
  </ThemeProvider>
);

export default Header;
