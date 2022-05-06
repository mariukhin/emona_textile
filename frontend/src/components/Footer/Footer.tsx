// node modules
import React from 'react';
// modules
import { colors } from 'utils/color';
// components
import ContactsBlock from './components/ContactsBlock';
import FooterInfoBlock from './components/FooterInfoBlock';
import { IconButton } from '@mui/material';
import { LinkedIn, Telegram } from '@mui/icons-material';
// styles
import {
  FooterWrapper,
  LogoContainer,
  LogoDescription,
  InfoWrapper,
  StyledDivider,
  DevelopersInfoWrapper,
  AllRightsReserved,
  DevelopersInfoContainer,
  DevelopersInfoBlock,
  DeveloperDetails,
  SocialButtonsBlock,
} from './styles';
import {
  BlockContainer,
  BlockHeading,
} from './components/FooterInfoBlock/styles';

const firstBlock = {
  title: 'Сторінки',
  subItems: ['Головна', 'Каталог', 'Про нас', 'Контакти'],
};

const secondBlock = {
  title: 'Каталог',
  subItems: [
    'Постільна білизна',
    'Постільні пріналежності',
    'Столова білизна',
    'Махрові вироби',
    'Рекламно-сувенірне',
    'Стрейчові чохли',
  ],
};

const Footer = () => (
  <FooterWrapper>
    <InfoWrapper>
      <LogoContainer>
        <img src="assets/logo-white.svg" alt="Emona logo" />
        <LogoDescription>
          Ми такі-то, займаємось тим-то та раді співпрацювати з вами
        </LogoDescription>
      </LogoContainer>

      <FooterInfoBlock
        title={firstBlock.title}
        subItems={firstBlock.subItems}
      />
      <FooterInfoBlock
        title={secondBlock.title}
        subItems={secondBlock.subItems}
      />

      <BlockContainer>
        <BlockHeading>Контакти</BlockHeading>
        <ContactsBlock />
      </BlockContainer>
    </InfoWrapper>

    <StyledDivider />

    <DevelopersInfoWrapper>
      <AllRightsReserved>
        © 2022 - ООО "Эмона Текстиль". Всі права захищені
      </AllRightsReserved>

      <DevelopersInfoContainer>
        <DevelopersInfoBlock>
          <DeveloperDetails>
            Дизайн — Костя Петруша
          </DeveloperDetails>
          <SocialButtonsBlock>
            <IconButton size="small" href='https://t.me/kostya_pet' target="_blank">
              <Telegram sx={{ color: colors.background.white }} />
            </IconButton>
            <IconButton size="small" href="https://www.linkedin.com/in/kostya-petrusha-a06355144/" target="_blank">
              <LinkedIn sx={{ color: colors.background.white }} />
            </IconButton>
          </SocialButtonsBlock>
        </DevelopersInfoBlock>
        <DevelopersInfoBlock>
          <DeveloperDetails>
            Розробка — Максим Марюхін
          </DeveloperDetails>
          <SocialButtonsBlock>
            <IconButton size="small" href='https://t.me/mar_max13' target="_blank">
              <Telegram sx={{ color: colors.background.white }} />
            </IconButton>
            <IconButton size="small" href='https://www.linkedin.com/in/maksym-mariukhin/' target="_blank">
              <LinkedIn sx={{ color: colors.background.white }} />
            </IconButton>
          </SocialButtonsBlock>
        </DevelopersInfoBlock>
      </DevelopersInfoContainer>
    </DevelopersInfoWrapper>
  </FooterWrapper>
);

export default Footer;
