// node modules
import React from 'react';
// components
import Header from 'components/Header';
import ScrollToTop from 'components/ScrollToTop';
import {
  Toolbar,
} from "@mui/material";
// styles
import styled from 'styled-components';

interface LayoutProps {
  children: React.ReactNode;
}

const ChildrenWrapper = styled.div`
  padding: 20px 0;
`;

const AppLayout: React.FC<LayoutProps> = ({ children }) => (
  <React.Fragment>
    <Header />
    <Toolbar id="back-to-top-anchor" />

    <ChildrenWrapper>
      {children}
    </ChildrenWrapper>

    <ScrollToTop />
  </React.Fragment>
);

export default AppLayout;
