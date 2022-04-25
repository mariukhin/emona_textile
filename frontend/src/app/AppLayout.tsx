// node modules
import React from 'react';
// components
import Header from 'components/Header';

interface LayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FC<LayoutProps> = ({ children }) => (
  <React.Fragment>
    <Header />

    {children}
  </React.Fragment>
);

export default AppLayout;
