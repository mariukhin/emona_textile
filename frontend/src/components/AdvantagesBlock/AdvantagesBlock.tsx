// node modules
import React from 'react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import AdvantagesBlockItem from './components/AdvantagesBlockItem';
// styles
import { AdvantagesBlockWrapper, StyledGridContainer } from './styles';
import { StyledButtonText, StyledButtonWrapper, StyledButton } from 'components/Header/styles';
import { colors } from 'utils/color';
import { goToForm } from 'modules';

type AdvantagesBlockProps = {
  advantageItems: AdvantagesBlockData[];
};

const AdvantagesBlock: React.FC<AdvantagesBlockProps> = ({
  advantageItems,
}) => (
  <AdvantagesBlockWrapper>
    <BlockInfoComponent title="Наші переваги" subtitle="Чому саме EMONA" />
  
    <StyledGridContainer>
      {advantageItems.map((item) => (
        <AdvantagesBlockItem
          key={item.id}
          title={item.title}
          subtitle={item.subtitle}
          iconUrl={item.iconUrl}
        />
      ))}
    </StyledGridContainer>
  
    <StyledButtonWrapper>
      <StyledButton onClick={ () => goToForm()} color="success" variant="contained" size="small">
        <StyledButtonText color={colors.text.white}>
          Зв’язатися з нами
        </StyledButtonText>
      </StyledButton>
    </StyledButtonWrapper>
  </AdvantagesBlockWrapper>
  ) 

export default AdvantagesBlock;
