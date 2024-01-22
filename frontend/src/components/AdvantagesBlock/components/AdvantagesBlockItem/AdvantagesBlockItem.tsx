// node modules
import React from 'react';
// styles
import {
  StyledPaper,
  BlockHeading,
  BlockHeadingContainer,
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
      <BlockHeadingContainer>
        <BlockHeading>{title}</BlockHeading>
      </BlockHeadingContainer>
      <BlockSubHeading>{subtitle}</BlockSubHeading>
    </BlockWrapper>
  </StyledPaper>
);

export default AdvantagesBlockItem;
