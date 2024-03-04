// node modules
import React from 'react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import AdvantagesBlockItem from './components/AdvantagesBlockItem';
// styles
import { AdvantagesBlockWrapper, StyledGridContainer, StyledButton } from './styles';
import { StyledButtonText, StyledButtonWrapper } from 'components/Header/styles';
// import ContactUsButton from 'components/ContactUsButton';
import { colors } from 'utils/color';

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
      <StyledButton color="success" variant="contained" size="small">
      {/* <ContactUsButton buttonElem={ () => <StyledButton color="success" variant="contained" size="small"/> }> */}
        <StyledButtonText color={colors.text.white}>
          Зв’язатися з нами
        </StyledButtonText>
      </StyledButton>
    </StyledButtonWrapper>
  </AdvantagesBlockWrapper>
);

export default AdvantagesBlock;
