// node modules
import React from 'react';
// modules
import { ROUTES } from 'routing/registration';
// import { Link as RouterLink } from "react-router-dom";
import styled from 'styled-components'
import { colors } from 'utils/color';
// components
import {
  AppBar,
  Button,
  Typography,
} from "@mui/material";

const headersData = [
  {
    label: "Головна",
    href: ROUTES.HOME,
  },
  {
    label: "Каталог",
    href: ROUTES.HOME,
  },
  {
    label: "Про нас",
    href: ROUTES.HOME,
  },
];

const StyledAppBar = styled(AppBar)`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  background-color: ${colors.background.white};
`;


const Header = () => (
  <header>
    <StyledAppBar>
      {headersData.map(({ label, href }) => (
          <Button
            {...{
              key: label,
              color: "primary",
            }}
          >
            <Typography variant='body1' textTransform="none">{label}</Typography>
          </Button>
        )
      )}
    </StyledAppBar>
  </header>
);

export default Header;
