// node modules
import React from 'react';
// modules
import { colors } from 'utils/color';
import { ROUTES } from 'routing/registration';
// components
import ContactsBlock from './components/ContactsBlock';
import FooterInfoBlock from './components/FooterInfoBlock';
import { IconButton } from '@mui/material';
import { LinkedIn, Telegram } from '@mui/icons-material';
// styles
import {
  FooterWrapper,
  Logo,
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
import { goToForm } from 'modules';

const firstBlock: FooterData = {
  title: 'Сторінки',
  subItems: [
    { label: 'Головна', href: ROUTES.HOME },
    { label: 'Каталог', href: ROUTES.CATALOG },
    { label: 'Про нас', href: ROUTES.ABOUT },
    { label: 'Контакти', onClick: () => goToForm() },
  ],
};

const secondBlock: FooterData = {
  title: 'Каталог',
  subItems: [
    { label: 'Постільна білизна', href: `${ROUTES.CATALOG_ITEM}?title=Постільна білизна` },
    { label: 'Столова білизна', href: `${ROUTES.CATALOG_ITEM}?title=Столова білизна` },
    { label: 'Махрові вироби', href: `${ROUTES.CATALOG_ITEM}?title=Махрові вироби` },
    { label: 'SPA & Resorts', href: `${ROUTES.CATALOG_ITEM}?title=SPA & Resorts` },
    { label: 'Рекламно-сувенірні вироби', href: `${ROUTES.CATALOG_ITEM}?title=Рекламно-сувенірні вироби` },
    { label: 'Чохли на меблі та декоративні вироби', href: `${ROUTES.CATALOG_ITEM}?title=Чохли на меблі та декоративні вироби` },
  ],
};

const Footer = () => (
  <FooterWrapper>
    <InfoWrapper>
      <LogoContainer>
        <Logo src="assets/logo-white.svg" alt="Emona logo" />
        <LogoDescription>
          Emona textile - компанія з продажу тканин та пошиття текстильних виробів для підприємств у сфері HoReCa. Надійність - це головна перевага нашої компанії.
        </LogoDescription>
      </LogoContainer>

      <FooterInfoBlock
        title={firstBlock.title}
        subItems={firstBlock.subItems}
        isFooter
      />
      <FooterInfoBlock
        title={secondBlock.title}
        subItems={secondBlock.subItems}
        isFooter
      />

      <BlockContainer isFooter title={'Контакти'}>
        <BlockHeading>Контакти</BlockHeading>
        <ContactsBlock isFooter />
      </BlockContainer>
    </InfoWrapper>

    <StyledDivider />

    <DevelopersInfoWrapper>
      <AllRightsReserved>
        © 2024 - ООО "Емона Текстиль". Всі права захищені
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
