import styled from 'styled-components';
import { colors } from 'utils/color';
import { Typography, Link } from '@mui/material';

export const BlockContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
`;

export const BlockHeading = styled(Typography)`
  font-size: 18px;
  line-height: 28px;
  font-weight: 700;
  color: ${colors.text.white};
  margin-bottom: 16px;
  font-family: 'Comfortaa';
`;

export const BlockSubItem = styled(Link)`
  font-size: 16px;
  line-height: 22px;
  color: ${colors.text.grey};
  margin-bottom: 8px;
  text-decoration: none;
  cursor: pointer;
  font-family: 'Montserrat';

  :hover {
    color: ${colors.text.orange};
  }
`;

export const ContactItem = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  margin-bottom: 8px;
`;

export const ContactBlockTextLink = styled(Link)`
  font-size: 16px;
  line-height: 22px;
  margin-left: 15px;
  color: ${colors.text.grey};
  font-family: 'Montserrat';
  text-decoration: none;

  :hover {
    text-decoration: underline;
  }
`;

export const CallIconWrapper = styled.div`
  margin-top: 0;
`;

export const ContactItemPhoneBlock = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
`;

