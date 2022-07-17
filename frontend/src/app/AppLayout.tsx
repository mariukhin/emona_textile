// node modules
import React from 'react';
import { observer } from 'mobx-react';
// components
import Header from 'components/Header';
import Footer from 'components/Footer';
import ScrollToTop from 'components/ScrollToTop';
import {
  CssBaseline,
  createTheme,
  ThemeProvider,
  Toolbar,
} from "@mui/material";
// styles
import styled from 'styled-components';

interface LayoutProps {
  children: React.ReactNode;
}

const ChildrenWrapper = styled.div`
  overflow: hidden;
  width: 100%;
`;

const theme = createTheme({
  breakpoints: {
    values: {
      xs: 320, // phone
      sm: 768, // tablets
      md: 1024, // small laptop
      lg: 1440, // desktop
      xl: 2560 // large screens
    }
  },
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

const AppLayout: React.FC<LayoutProps> = ({ children }) => (
  <ThemeProvider theme={theme}>
    <CssBaseline />
    
    <Header />
    <Toolbar id="back-to-top-anchor" />

    <ChildrenWrapper>
      {children}
    </ChildrenWrapper>

    <ScrollToTop />
    <Footer />
  </ThemeProvider>
);

export default observer(AppLayout);
