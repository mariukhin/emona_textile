import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';
import { Typography, Link } from '@mui/material';

export const ContactsBlockContainer = styled.div<{ isFooter: boolean; }>`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
`;

export const BlockContainer = styled.div<{ isFooter: boolean; title: string }>`
  flex-direction: column;
  justify-content: flex-start;

  @media ${device.mobile} {
    display: ${props => props.isFooter ? 'none' : 'flex'};
    margin-bottom: ${props => props.title === 'Каталог'  ? '32px' : '60px'};
  }

  @media ${device.tablet} {
    display: ${props => props.title === 'Сторінки' || props.title === 'Каталог'  ? 'none' : 'flex'};
    max-width: ${props => props.isFooter && '280px'};
    margin-bottom: 0;
  }

  @media ${device.laptop} {
    display: flex;
    max-width: ${props => props.isFooter && props.title === 'Каталог' && '220px'};
  }
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
  margin-bottom: 12px;
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

