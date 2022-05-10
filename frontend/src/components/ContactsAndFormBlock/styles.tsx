import styled from 'styled-components';
import { colors } from 'utils/color';
import { Paper, Button, Typography, FormControl } from '@mui/material';

export const ContactsAndFormBlockWrapper = styled.div`
  margin: 0 0 80px;
`;

export const BlockContainer = styled.div`
  display: flex;
  justify-content: space-between;
  width: 75%;
  margin: 0 auto;
  padding-top: 40px;
`;

export const InfoContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 40%;
  height: 430px;
  padding: 30px 0 0;
`;

export const StyledMapComponent = styled.div`
  border: 2px solid ${colors.background.green};
  border-radius: 20px;
  width: 100%;
  height: 248px;
`;

export const StyledPaper = styled(Paper)`
  width: 536px;
  border-radius: 20px;
  box-shadow: 0px 6px 24px rgba(0, 0, 0, 0.12);
  background-color: #fff;
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

export const StyledFormControl = styled(FormControl)`
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
`;

export const StyledButtonText = styled(Typography)`
  font-size: 16px;
  line-height: 26px;
  text-transform: uppercase;
  font-family: 'Nunito';
  font-weight: 700;
`;
