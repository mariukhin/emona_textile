// node modules
import React from 'react';
// components
import Catalog from 'components/Catalog';
// mocks
import { mockedCatalogItems } from 'components/Catalog/mocks';
// styles
import {
  CatalogPhotoContainer,
  CarouselHeading,
  InfoBlock,
  StyledButton,
  StyledButtonText,
} from './styles';

const CatalogScreenView = () => {
  return (
    <div>
      <CatalogPhotoContainer>
        <InfoBlock>
          <CarouselHeading>Каталог</CarouselHeading>
          <StyledButton color="warning" size="large" variant="contained">
            <StyledButtonText>
              Якісні товари
            </StyledButtonText>
          </StyledButton>
        </InfoBlock>
      </CatalogPhotoContainer>

      <Catalog catalogItems={mockedCatalogItems} />
    </div>
  );
};

export default CatalogScreenView;
