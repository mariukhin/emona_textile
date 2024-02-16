// node modules
import React from 'react';
// modules
import { colors } from 'utils/color';
// components
import {
  Call,
  LocationOn,
  EmailOutlined,
} from '@mui/icons-material';
// styles
import {
  ContactsBlockContainer,
  ContactItem,
  ContactItemPhoneBlock,
  ContactBlockTextLink,
} from '../FooterInfoBlock/styles';

type ContactsBlockType = {
  isFooter?: boolean;
};

const ContactsBlock: React.FC<ContactsBlockType> = ({ isFooter = false }) => {
  const iconColor = isFooter ? colors.background.white : colors.background.green;
  const textColor = isFooter ? colors.text.grey : colors.text.greyDark;

  return (
    <ContactsBlockContainer isFooter={ isFooter }>
      <ContactItem>
        <LocationOn sx={{ color: iconColor }} />
        <ContactBlockTextLink
          href="https://maps.google.com?q=50.454905968545766, 30.48833512571684"
          target="_blank"
          sx={{ color: textColor }}
        >
          Україна, 01135, Київ, вул. Дмитрівська 82, офіс 86
        </ContactBlockTextLink>
      </ContactItem>
      <ContactItem>
        <Call sx={{ color: iconColor }} />
        <ContactItemPhoneBlock>
          <ContactBlockTextLink href="tel:+380444868610" sx={{ color: textColor }}>
            +38 044 486 86 10
          </ContactBlockTextLink>
          <ContactBlockTextLink href="tel:+380444868596" sx={{ color: textColor }}>
            +38 044 486 85 96
          </ContactBlockTextLink>
        </ContactItemPhoneBlock>
      </ContactItem>
      <ContactItem>
        <EmailOutlined sx={{ color: iconColor }} />
        <ContactBlockTextLink
          href="mailto:emona.textile@gmail.com?subject=To Emona Textile company"
          sx={{ color: textColor }}
        >
          emona.textile@gmail.com
        </ContactBlockTextLink>
      </ContactItem>
    </ContactsBlockContainer>
  );
};

export default ContactsBlock;
