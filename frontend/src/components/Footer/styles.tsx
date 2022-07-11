import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Divider, Typography } from '@mui/material';

export const FooterWrapper = styled.footer`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 100%;
  background-color: ${colors.background.green};

  @media ${device.mobile} {
    padding: 60px 0 0;
  }

  @media ${device.tablet} {
    padding: 60px 24px 0;
  }

  @media ${device.laptop} {
    padding: 60px 64px 0;
  }
`;

export const InfoWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;

  @media ${device.mobile} {
    flex-direction: column;
    padding-left: 77px;
    padding-right: 77px;
  }

  @media ${device.tablet} {
    flex-direction: row;
    padding: 0 0 64px;
  }
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

  @media ${device.mobile} {
    align-content: center;
    text-align: center;
    width: 100%;
    min-height: 115px;
    max-height: 130px;
    margin-bottom: 50px;
  }

  @media ${device.tablet} {
    align-items: flex-start;
    text-align: left;
    width: 35%;
    max-width: 250px;
    min-height: 90px;
    max-height: 100px;
    margin-bottom: 0;
  }

  @media ${device.laptop} {
    width: 25%;
    max-width: 250px;
    min-height: 80px;
    max-height: 95px;
  }

  @media ${device.laptopL} {
    width: 20%;
    max-width: 250px;
    min-height: 100px;
    max-height: 110px;
  }
`;

export const StyledDivider = styled(Divider)`
  border: 1px solid #4CAF50;

  @media ${device.mobile} {
    margin-left: 60px;
    margin-right: 60px;
  }

  @media ${device.tablet} {
    margin: 0;
  }
`;

export const DevelopersInfoWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  text-align: center;

  @media ${device.mobile} {
    margin: 20px auto;
    width: 339px;
  }

  @media ${device.tablet} {
    width: 530px;
    margin: 32px auto;
  }
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

  @media ${device.mobile} {
    flex-direction: column;
    align-items: center;
  }

  @media ${device.tablet} {
    flex-direction: row;
  }
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
