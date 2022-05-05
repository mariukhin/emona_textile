// node modules
import React from 'react';
// modules
import { colors } from 'utils/color';
// components
import {
  Typography,
} from '@mui/material';
// styles
import {
  FooterWrapper,
  Logo,
  LogoContainer,
  InfoWrapper,
  StyledDivider,
} from './styles';

const Footer = () => (
  <div style={{ height: '468px', width: '100%', backgroundColor: colors.background.green }}>
    <FooterWrapper>
      <InfoWrapper>
        <LogoContainer>
          <Logo src='assets/logo-white.svg' alt="Emona logo" />
    
          <Typography
            variant="body1"
            fontSize="14px"
            fontWeight={400}
            color={colors.text.white}
            sx={{
              fontFamily: 'Montserrat',
            }}
          >Ми такі-то, займаємось тим-то та раді співпрацювати з вами</Typography>
        </LogoContainer>
      </InfoWrapper>

      <StyledDivider />
    </FooterWrapper>
  </div>
);

export default Footer;
