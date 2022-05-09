// node modules
import React from 'react';
// modules
import { colors } from 'utils/color';
// components
import {
  CallOutlined,
  LocationOnOutlined,
  EmailOutlined,
} from '@mui/icons-material';
// styles
import {
  BlockContainer,
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
    <BlockContainer>
      <ContactItem>
        <LocationOnOutlined sx={{ color: iconColor }} />
        <ContactBlockTextLink
          href="https://maps.google.com?q=50.454905968545766, 30.48833512571684"
          target="_blank"
          sx={{ color: textColor }}
        >
          Київ, вул. Дмитрівська 82, офіс 1
        </ContactBlockTextLink>
      </ContactItem>
      <ContactItem>
        <CallOutlined sx={{ color: iconColor }} />
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
    </BlockContainer>
  );
};

export default ContactsBlock;
