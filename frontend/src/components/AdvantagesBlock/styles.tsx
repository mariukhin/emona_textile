import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Paper, Button, Typography } from '@mui/material';

export const AdvantagesBlockWrapper = styled.div`
  margin: 120px 0;
`;

export const StyledGridContainer = styled.div`
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;

  @media ${device.mobile} {
    width: 92%;
    margin-top: 58px;
  }

  @media ${device.tablet} {
    width: 95%;
    margin-top: 40px;
  }

  @media ${device.laptopL} {
    width: 90%;
  }
`;

export const StyledPaper = styled(Paper)`
  border-radius: 20px;

  @media ${device.mobile} {
    width: 100%;
    margin-bottom: 23px;

    :last-child {
      margin-bottom: 0;
    }
  }

  @media ${device.tablet} {
    width: 48.5%;

    :nth-last-child(2)  {
      margin-bottom: 0;
    }
  }

  @media ${device.laptopL} {
    width: 24%;
    margin-bottom: 0;
  }
`;

export const BlockWrapper = styled.div`
  width: 100%;
  padding: 20px;
`;

export const BlockImage = styled.img`
  width: 80px;
  height: 80px;
  margin-bottom: 25px;
`;

export const BlockHeadingContainer = styled.div`
  @media ${device.mobile} {
    margin-bottom: 10px;
  }

  @media ${device.laptopL} {
    height: 72px;
  }
`;

export const BlockHeading = styled(Typography)`
  font-size: 24px;
  line-height: 30px;
  color: ${colors.text.greyDark};
  font-family: 'Comfortaa';
  font-weight: 700;
`;

export const BlockSubHeading = styled(Typography)`
  font-size: 16px;
  line-height: 24px;
  color: ${colors.text.greyLight};
  font-family: 'Montserrat';
  font-weight: 400;
  letter-spacing: 0.48px;
`;

export const StyledButton = styled(Button)`
  padding: 10px 28px;
  background-color: ${colors.button.default};
`;
