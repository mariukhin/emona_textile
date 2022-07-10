import styled from 'styled-components';
import { colors } from 'utils/color';
import { Divider, Typography } from '@mui/material';

export const FooterWrapper = styled.footer`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 60px 64px 0;
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

export const LogoDescription = styled(Typography)`
  font-size: 14px;
  font-weight: 400;
  color: ${colors.text.white};
  font-family: 'Montserrat';
`;

export const LogoContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 20%;
  max-width: 250px;
  min-height: 100px;
  max-height: 110px;
`;

export const StyledDivider = styled(Divider)`
  border: 1px solid #4CAF50;
`;

export const DevelopersInfoWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  text-align: center;
  margin: 32px auto;
  width: 530px;
`;

export const AllRightsReserved = styled(Typography)`
  font-size: 12px;
  font-weight: 400;
  color: ${colors.text.grey};
  font-family: 'Montserrat';
`;

export const DevelopersInfoContainer = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  margin-top: 18px;
`;

export const DevelopersInfoBlock = styled.div`
  display: flex;
  justify-content: flex-start;
  align-items: center;
`;

export const SocialButtonsBlock = styled.div`
  display: flex;
  justify-content: space-between;
`;

export const DeveloperDetails = styled(Typography)`
  font-size: 12px;
  font-weight: 400;
  color: ${colors.text.white};
  font-family: 'Montserrat';
`;
