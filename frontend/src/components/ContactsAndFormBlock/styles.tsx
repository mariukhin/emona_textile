import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Paper, Button, Typography, FormControl } from '@mui/material';

export const ContactsAndFormBlockWrapper = styled.div`
  margin: 0 0 80px;
`;

export const BlockContainer = styled.div`
  display: flex;
  margin: 0 auto;
  padding-top: 40px;

  @media ${device.mobile} {
    flex-direction: column-reverse;
    align-items: center;
    width: 95%;
  }

  @media ${device.laptopL} {
    justify-content: space-between;
    flex-direction: row;
    width: 75%;
  }
`;

export const InfoContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;

  @media ${device.mobile} {
    padding: 60px 0 0;
    width: 100%;
    height: 445px;
    align-items: center;
  }

  @media ${device.tablet} {
    height: 630px;
  }

  @media ${device.laptop} {
    height: 780px;
  }

  @media ${device.laptopL} {
    align-items: flex-start;
    padding: 30px 0 0;
    width: 40%;
    height: 430px;
  }
`;

export const StyledMapComponent = styled.div`
  border: 2px solid ${colors.background.green};
  border-radius: 20px;
  width: 100%;

  @media ${device.mobile} {
    height: 224px;
  }

  @media ${device.tablet} {
    height: 418px;
  }

  @media ${device.laptop} {
    height: 567px;
  }

  @media ${device.laptopL} {
    height: 248px;
  }
`;

export const StyledPaper = styled(Paper)<{ isErrors: boolean; }>`
  border-radius: 20px;
  box-shadow: 0px 6px 24px rgba(0, 0, 0, 0.12);
  background-color: #fff;

  @media ${device.mobile} {
    width: 100%;
    height: ${props => props.isErrors ? '650px' : '564px'};
  }

  @media ${device.tablet} {
    width: 536px;
    height: ${props => props.isErrors ? '550px' : '484px'};
  }
`;

export const PaperWrapper = styled.div`
  padding: 39px;
`;

export const FormHeader = styled(Typography)`
  font-size: 24px;
  line-height: 30px;
  color: ${colors.text.greyDark};
  margin: 0 0;
  font-family: 'Comfortaa';
  font-weight: 700;
  text-align: center;
`;

export const FormBlock = styled.div`
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
  width: 100%;

  @media ${device.mobile} {
    flex-direction: column;
  }

  @media ${device.tablet} {
    flex-direction: row;
  }
`;

export const StyledFormControl = styled(FormControl)<{ width: string; marginCustom?: string; }>`
  @media ${device.mobile} {
    width: 100%;

    :not(:first-of-type) {
      margin: 12px 0 0;
    }
  }

  @media ${device.tablet} {
    width: ${ props => props.width};
    margin: ${ props => props.marginCustom || 0} !important;
  }

  .Mui-focused {
    color: ${colors.background.green};

    .MuiOutlinedInput-notchedOutline {
      border-color: ${colors.background.green} !important;
    }
  }
`;

export const StyledButton = styled(Button)`
  padding: 10px 20px 10px 28px;
  width: 100%;

  @media ${device.mobile} {
    margin-top: 12px;
  }

  @media ${device.tablet} {
    margin-top: 0;
  }
`;

export const StyledButtonText = styled(Typography)`
  font-size: 16px;
  line-height: 26px;
  text-transform: uppercase;
  font-family: 'Nunito';
  font-weight: 700;
`;
