// node modules
import React from 'react';
// components
import PagePhotoBlock from 'components/PagePhotoBlock';
import { mockedCatalogItems } from './mocks';
import { ArrowForward } from '@mui/icons-material';
// styles
import {
  StyledPaper,
  StyledGrid,
  ItemContainer,
  ItemImage,
  StyledGridContainer,
  ItemInfoBlock,
  ItemInfoBlockTitle,
  ItemButton,
  ItemButtonText,
  ItemInfoDescriptionList,
  ItemInfoDescriptionListItem,
  CatalogWrapper
} from './styles';

const CatalogItemScreenView = () => {
  const currentItem = mockedCatalogItems[0];
  return (
    <div>
      <PagePhotoBlock heading={ currentItem.title } btnText={ currentItem.buttonText } imageUrl={ currentItem.imageUrl } />
      <CatalogWrapper isMainPage>
        <StyledGridContainer container spacing={3}>
          {currentItem.items.map((item) => (
            <StyledGrid key={ item.id } item xs={12} sm={6} md={6} lg={6}>
              <StyledPaper>
                <ItemContainer>
                  <ItemImage
                    theme={{
                      main: item.imageUrl,
                    }}
                  />
                  <ItemInfoBlock>
                    <ItemInfoBlockTitle>{ item.title }</ItemInfoBlockTitle>
                    <ItemInfoDescriptionList>
                      { item.description.map(descr => <ItemInfoDescriptionListItem>{ descr };</ItemInfoDescriptionListItem>) }
                    </ItemInfoDescriptionList>
                    <ItemButton color="success" variant="text" endIcon={<ArrowForward />}>
                      <ItemButtonText>
                        ЗАМОВИТИ
                      </ItemButtonText>
                    </ItemButton>
                  </ItemInfoBlock>
                </ItemContainer>
              </StyledPaper>
            </StyledGrid>
          ))}
        </StyledGridContainer>
      </CatalogWrapper>
    </div>
  );
};

export default CatalogItemScreenView;
