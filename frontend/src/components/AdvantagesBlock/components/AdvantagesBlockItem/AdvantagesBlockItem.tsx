// node modules
import React from 'react';
// styles
import {
  StyledPaper,
  BlockHeading,
  BlockSubHeading,
  BlockWrapper,
  BlockImage,
} from '../../styles';

type AdvantagesBlockItemProps = {
  title: AdvantagesBlockData['title'];
  iconUrl: AdvantagesBlockData['iconUrl'];
  subtitle: AdvantagesBlockData['subtitle'];
};

const AdvantagesBlockItem: React.FC<AdvantagesBlockItemProps> = ({
  title,
  iconUrl,
  subtitle,
}) => (
  <StyledPaper elevation={0}>
    <BlockWrapper>
      <BlockImage src={iconUrl} alt="AdvantageBlockIcon" />
      <BlockHeading>{title}</BlockHeading>
      <BlockSubHeading>{subtitle}</BlockSubHeading>
    </BlockWrapper>
  </StyledPaper>
);

export default AdvantagesBlockItem;
