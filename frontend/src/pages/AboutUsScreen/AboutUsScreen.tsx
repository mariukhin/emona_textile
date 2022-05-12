// node modules
import React from 'react';
// components
import Banner from 'components/Banner';
import PagePhotoBlock from 'components/PagePhotoBlock';
// styles
import { AboutUsWrapper, AboutUsText, BannerContainer, BlockImage } from './styles';

const AboutUsScreenView = () => (
  <div>
    <PagePhotoBlock
      heading="Про нас"
      btnText="Досвід компанії"
      imageUrl="assets/about-us.jpeg"
    />

    <AboutUsWrapper>
      <AboutUsText>
        Як відомо, гостинність дуже велике і важливе поняття. Його важливість
        полягає у задоволенні всіх можливих бажань гостя. Вибираючи готель,
        мандрівник припускає не лише відпочити і виспатися, а й привести себе до
        ладу в чужому місті. І не важливо – чи відпочиває людина у дорогому
        готелі на модному курорті, зупиняється у готелі економ класу в
        насиченому екскурсійному турі або з ранку до ночі зайнята переговорами з
        іноземними партнерами – при поверненні до готелю йому потрібен
        повноцінний комфортний відпочинок.
      </AboutUsText>

      <BannerContainer>
        <Banner title="96% клієнтів" subtitle="стають постійними" />
        <Banner title="15 років" subtitle="досвіду" />
        <Banner title="800 замовлень" subtitle="виготовлено" />
      </BannerContainer>

      <AboutUsText marginBottom="40px">
        Обов'язковим атрибутом кожного готелю та ресторану є текстиль. Це і
        постільні приналежності, і постільна білизна, і махрові вироби, і
        столова білизна, штори, покривала та ін.
      </AboutUsText>
      <AboutUsText marginBottom="40px">
        Фірма «Емона Текстиль» надає послуги з комплектації готелів та
        ресторанів професійними текстильними виробами.
      </AboutUsText>

      <BlockImage src="assets/CataloguePhotos/Постільна-білизна.jpeg"/>

      <AboutUsText marginTop="60px" marginBottom="40px">
        У нашому асортименті великий вибір готельно-ресторанного та домашнього текстилю:<br/>
        <br/>
        - для готелів - постільна білизна, постільні речі, махрові вироби<br/>
        - для ресторанів - скатертини, серветки, наперони, сети, фуршетні спідниці, чохли на стільці та меблі<br/>
        - для дому - постільна білизна, постільні речі, махрові вироби, скатертини, серветки, чохли на меблі, декоративні подушки, сувенірна текстильна продукція.
      </AboutUsText>

      <BlockImage src="assets/CataloguePhotos/рекламно-сувенірна-продукція.jpeg"/>
    </AboutUsWrapper>
  </div>
);

export default AboutUsScreenView;
