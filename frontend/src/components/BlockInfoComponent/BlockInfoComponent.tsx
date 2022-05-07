// node modules
import React from 'react';
// styles
import {
  BlockContainer,
  BlockHeading,
  BlockSubHeading,
} from './styles';

type BlockInfoComponentProps = {
  title: string;
  subtitle: string;
}

const BlockInfoComponent: React.FC<BlockInfoComponentProps> = ({ title, subtitle }) => (
  <BlockContainer>
    <BlockHeading>{title}</BlockHeading>
    <BlockSubHeading>{subtitle}</BlockSubHeading>
  </BlockContainer>
);

export default BlockInfoComponent;
