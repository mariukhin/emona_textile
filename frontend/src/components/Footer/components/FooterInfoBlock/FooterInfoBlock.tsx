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
  isFooter: boolean;
}

const FooterInfoBlock: React.FC<FooterInfoBlockComponentProps> = ({ title, subItems, isFooter = false }) => (
  <BlockContainer isFooter={ isFooter } title={ title }>
    <BlockHeading>{title}</BlockHeading>

    {subItems.map(item => <BlockSubItem key={useId()} onClick={ item.onClick && item.onClick } href={item.href && item.href}>{item.label}</BlockSubItem> )}
  </BlockContainer>
);

export default FooterInfoBlock;
