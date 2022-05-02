// node modules
import React from 'react';
// components
import {
  createTheme,
  ThemeProvider,
} from '@mui/material';
// styles
import {
  BlockContainer,
  BlockHeading,
  BlockSubHeading,
} from './styles';

type BlockInfoComponentProps = {
  title: string;
  subtitle: string;
}

const HeadingTheme = createTheme({
  typography: {
    fontFamily: 'Comfortaa',
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: `
        @font-face {
          font-family: 'Comfortaa';
          src: url('fonts/Comfortaa/Comfortaa-Bold.ttf');
          font-weight: 700;
        }
      `,
    },
  },
});

const SubHeadingTheme = createTheme({
  typography: {
    fontFamily: 'Montserrat',
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: `
        @font-face {
          font-family: 'Montserrat';
          src: url('fonts/Montserrat/Montserrat-SemiBold.ttf');
          font-weight: 600;
        }
      `,
    },
  },
});

const BlockInfoComponent: React.FC<BlockInfoComponentProps> = ({ title, subtitle }) => (
  <BlockContainer>
    <ThemeProvider theme={HeadingTheme}>
      <BlockHeading>{title}</BlockHeading>
    </ThemeProvider>
    <ThemeProvider theme={SubHeadingTheme}>
      <BlockSubHeading>{subtitle}</BlockSubHeading>
    </ThemeProvider>
  </BlockContainer>
);

export default BlockInfoComponent;
