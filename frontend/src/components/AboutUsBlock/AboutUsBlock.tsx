// node modules
import React from 'react';
import { useStore } from 'modules/Stores';
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
} from './styles';
import { colors } from 'utils/color';

const AboutUsBlock = () => {
  const { goToRoute } = useStore('RoutingStore');

  const handleClickDetails = () => goToRoute(ROUTES.ABOUT);

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
            Як відомо, гостинність дуже широке і важливе поняття. Його важливість полягає у задоволенні всіх можливих бажань гостя. Обираючи готель, мандрівник передбачає не лише відпочити і виспатися, а й привести себе до ладу в чужому місті. І не важливо чи відпочиває людина у дорогому готелі на модному курорті, зупиняється у готелі економ класу в насиченому екскурсійному турі або з ранку до ночі зайнята переговорами з іноземними партнерами при поверненні до готелю йому потрібен повноцінний комфортний відпочинок.
            Обов'язковим атрибутом кожного готелю та ресторану є текстиль. Це і постільні речі, і постільна білизна, і махрові вироби, і столова білизна, штори, покривала та ін.
            Фірма "Емона Текстиль" понад 20 років надає послуги з комплектації готелів та ресторанів професійними текстильними виробами. У нашому асортименті є постільні речі (подушки, ковдри, наматрацники), постільна білизна (наволочки, підковдри, простирадла), столова білизна (скатертини, серветки, сети, фуршетні спідниці, чохли для стільців і столів) і махрові вироби (рушники, халати). Ми здійснюємо прямі поставки текстильних виробів та тканин від провідних виробників України, Чехії, Португалії, Іспанії, Туреччини.
            Тісні ділові контакти з іноземними заводами-виробниками, знання українського текстильного ринку, великий досвід роботи та наявність власного швейного виробництва дозволяє оперативно та якісно виконувати будь-які замовлення наших клієнтів. Співпраця з текстильними заводами різних країн дає можливість поєднувати високу якість та доступні ціни.
          </BlockText>
  
          <BannerContainer>
            <Banner title="96% клієнтів" subtitle="стають постійними" />
            <Banner title="15 років" subtitle="досвіду" />
          </BannerContainer>
  
          <StyledButton
            color="warning"
            variant="contained"
            size="small"
            endIcon={<ArrowForward />}
            onClick={handleClickDetails}
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
