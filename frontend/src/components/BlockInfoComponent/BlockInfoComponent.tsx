// node modules
import React from 'react';
// styles
import {
  BlockHeading,
  BlockSubHeading,
} from './styles';

type BlockInfoComponentProps = {
  title: string;
  subtitle: string;
  textAlign?: CanvasTextAlign;
}

const BlockInfoComponent: React.FC<BlockInfoComponentProps> = ({ title, subtitle, textAlign = 'center' }) => (
  <div style={{ textAlign }}>
    <BlockHeading>{title}</BlockHeading>
    <BlockSubHeading>{subtitle}</BlockSubHeading>
  </div>
);

export default BlockInfoComponent;
