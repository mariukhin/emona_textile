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
            Як відомо, гостинність дуже велике і важливе поняття. Його важливість
            полягає у задоволенні всіх можливих бажань гостя. Вибираючи готель,
            мандрівник припускає не лише відпочити і виспатися, а й привести себе
            до ладу в чужому місті. І не важливо – чи відпочиває людина у дорогому
            готелі на модному курорті, зупиняється у готелі економ класу в
            насиченому екскурсійному турі або з ранку до ночі зайнята переговорами
            з іноземними партнерами – при поверненні до готелю йому потрібен
            повноцінний комфортний відпочинок.
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
