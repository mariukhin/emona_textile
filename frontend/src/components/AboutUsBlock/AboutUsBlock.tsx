// node modules
import React from 'react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import Banner from 'components/Banner';
import { ArrowForward } from '@mui/icons-material';
import { ROUTES } from 'routing/registration';
// styles
import {
  AboutUsBlockWrapper,
  AboutUsPhotoBlock,
  StyledPaper,
  PaperWrapper,
  BlockText,
  BannerContainer,
  StyledButton,
  StyledButtonText,
  BlockTextItem,
} from './styles';
import { colors } from 'utils/color';

const AboutUsBlock = () => {
  return (
    <AboutUsBlockWrapper>
      <AboutUsPhotoBlock />
  
      <StyledPaper>
        <PaperWrapper>
          <BlockInfoComponent
            title="Про нас"
            subtitle="Досвід компанії"
            textAlign="left"
          />
          <BlockText>
            <BlockTextItem>
              Як відомо, гостинність - дуже широке і важливе поняття. Його важливість полягає у задоволенні всіх можливих бажань гостя. Обов'язковим атрибутом кожного готелю та ресторану є текстиль.
            </BlockTextItem>

            <BlockTextItem>
              Фірма "Емона Текстиль" понад 20 років надає послуги з комплектації готелів та ресторанів професійними текстильними виробами. У нашому асортименті є постільна білизна, постільні речі,, столова білизна, фуршетні спідниці, чохли для стільців і столів, махрові вироби. Ми здійснюємо прямі поставки текстильних виробів та тканин від провідних виробників України, Чехії, Португалії, Іспанії, Туреччини. 
            </BlockTextItem>

            <BlockTextItem>
              Тісні ділові контакти з іноземними заводами-виробниками, знання українського текстильного ринку, великий досвід роботи та наявність власного швейного виробництва дозволяє оперативно та якісно виконувати будь-які замовлення наших клієнтів. 
            </BlockTextItem>
          </BlockText>
  
          <BannerContainer>
            <Banner title="96% клієнтів" subtitle="стають постійними" />
            <Banner title="20 років" subtitle="досвіду" />
          </BannerContainer>
  
          <StyledButton
            color="warning"
            variant="contained"
            size="small"
            endIcon={<ArrowForward />}
            href={ ROUTES.ABOUT }
          >
            <StyledButtonText variant="button" color={colors.text.white}>
              Детальніше
            </StyledButtonText>
          </StyledButton>
        </PaperWrapper>
      </StyledPaper>
    </AboutUsBlockWrapper>
  )
};

export default AboutUsBlock;
