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
    <BlockHeading sx={{ fontFamily: 'Comfortaa' }}>{title}</BlockHeading>
    <BlockSubHeading sx={{ fontFamily: 'Montserrat' }}>{subtitle}</BlockSubHeading>
  </BlockContainer>
);

export default BlockInfoComponent;
