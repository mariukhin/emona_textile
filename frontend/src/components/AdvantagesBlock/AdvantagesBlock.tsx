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
import { useAnalyticsEventTracker } from 'modules';

type AdvantagesBlockProps = {
  advantageItems: AdvantagesBlockData[];
};

const AdvantagesBlock: React.FC<AdvantagesBlockProps> = ({
  advantageItems,
}) => {
  const gaEventTracker = useAnalyticsEventTracker('Зв\'язатись з нами');

  const handleGoToForm = () => {
    gaEventTracker('Кнопка');
    goToForm();
  }

  return (
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
        <StyledButton onClick={handleGoToForm} color="success" variant="contained" size="small">
          <StyledButtonText color={colors.text.white}>
            Зв’язатися з нами
          </StyledButtonText>
        </StyledButton>
      </StyledButtonWrapper>
    </AdvantagesBlockWrapper>
  );
};

export default AdvantagesBlock;
