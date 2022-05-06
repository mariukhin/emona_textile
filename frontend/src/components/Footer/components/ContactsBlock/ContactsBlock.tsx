// node modules
import React from 'react';
// modules
import { colors } from 'utils/color';
// components
import { CallOutlined, LocationOnOutlined, EmailOutlined } from '@mui/icons-material';
// styles
import {
  BlockContainer,
  ContactItem,
  ContactItemPhoneBlock,
  ContactBlockTextLink,
} from '../FooterInfoBlock/styles';

const ContactsBlock = () => (
  <BlockContainer>
    <ContactItem>
      <LocationOnOutlined sx={{ color: colors.background.white }} />
      <ContactBlockTextLink
        href="https://maps.google.com?q=50.454905968545766, 30.48833512571684"
        target="_blank"
      >
        Київ, вул. Дмитрівська 82, офіс 1
      </ContactBlockTextLink>
    </ContactItem>
    <ContactItem>
      <CallOutlined sx={{ color: colors.background.white }} />
      <ContactItemPhoneBlock>
        <ContactBlockTextLink href="tel:+380444868610">+38 044 486 86 10</ContactBlockTextLink>
        <ContactBlockTextLink href="tel:+380444868596">+38 044 486 85 96</ContactBlockTextLink>
      </ContactItemPhoneBlock>
    </ContactItem>
    <ContactItem>
      <EmailOutlined sx={{ color: colors.background.white }} />
      <ContactBlockTextLink
        href="mailto:emona.textile@gmail.com?subject=To Emona Textile company"
      >
        emona.textile@gmail.com
      </ContactBlockTextLink>
    </ContactItem>
  </BlockContainer>
);

export default ContactsBlock;
