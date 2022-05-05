import styled from 'styled-components';
import { colors } from 'utils/color';
import { Divider } from '@mui/material';

export const FooterWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 60px 64px;
  width: 100%;
  background-color: ${colors.background.green};
`;

export const InfoWrapper = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 0 0 64px;
  width: 100%;
`;

export const Logo = styled.img`
  padding: 0;

  svg {
    fill: ${colors.background.white};
  }
`;

export const LogoContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 20%;
  height: 100px;
`;

export const StyledDivider = styled(Divider)`
  border: 1px solid #4CAF50;
`;
