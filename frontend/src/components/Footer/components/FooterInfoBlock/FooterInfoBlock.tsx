// node modules
import React from 'react';
// styles
import {
  BlockContainer,
  BlockHeading,
  BlockSubItem,
} from './styles';

type FooterInfoBlockComponentProps = {
  title: string;
  subItems: string[];
}

const FooterInfoBlock: React.FC<FooterInfoBlockComponentProps> = ({ title, subItems }) => (
  <BlockContainer>
    <BlockHeading>{title}</BlockHeading>

    {subItems.map(item => <BlockSubItem>{item}</BlockSubItem> )}
  </BlockContainer>
);

export default FooterInfoBlock;
