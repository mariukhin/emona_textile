// node modules
import React from 'react';
// modules
import { ROUTES } from 'routing/registration';
import logo from 'assets/logo.svg';
// import { Link as RouterLink } from "react-router-dom";
import styled from 'styled-components'
import { colors } from 'utils/color';
// components
import {
  AppBar,
  Button,
  Typography,
  IconButton,
  Stack,
} from "@mui/material";

const headersData: HeadersData[] = [
  {
    label: "Головна",
    href: ROUTES.HOME,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: "Каталог",
    href: ROUTES.HOME,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: "Про нас",
    href: ROUTES.HOME,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: "Зв’язатися",
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
`;


const Header = () => (
  <header>
    <StyledAppBar>
      <IconButton>
        <img src={logo} alt="Emona logo" />
      </IconButton>

      <Stack direction="row" spacing={4}>
        {headersData.map(({ label, href, variant, color }) => (
            <Button
              {...{
                key: label,
                color: "success",
                href,
                variant,
                size: "small"
              }}
            >
              <Typography
                variant='body1'
                color={color}
                textTransform={variant === 'contained' ? 'uppercase' : 'none'}
              >
                {label}
              </Typography>
            </Button>
          )
        )}
      </Stack>
    </StyledAppBar>
  </header>
);

export default Header;
