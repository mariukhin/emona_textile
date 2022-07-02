// node modules
import React from 'react';
import { observer } from 'mobx-react';
// components
import Header from 'components/Header';
import Footer from 'components/Footer';
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
  width: 100%;
`;

const AppLayout: React.FC<LayoutProps> = ({ children }) => (
  <React.Fragment>
    <Header />
    <Toolbar id="back-to-top-anchor" />

    <ChildrenWrapper>
      {children}
    </ChildrenWrapper>

    <ScrollToTop />
    <Footer />
  </React.Fragment>
);

export default observer(AppLayout);
