// node modules
import React, { useId } from 'react';
// styles
import {
  BlockContainer,
  BlockHeading,
  BlockSubItem,
} from './styles';

type FooterInfoBlockComponentProps = {
  title: FooterData['title'];
  subItems: FooterData['subItems'];
}

const FooterInfoBlock: React.FC<FooterInfoBlockComponentProps> = ({ title, subItems }) => (
  <BlockContainer isFooter={false} title={ title }>
    <BlockHeading>{title}</BlockHeading>

    {subItems.map(item => <BlockSubItem key={useId()} href={item.href}>{item.label}</BlockSubItem> )}
  </BlockContainer>
);

export default FooterInfoBlock;
