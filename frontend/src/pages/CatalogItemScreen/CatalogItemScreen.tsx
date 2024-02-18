// node modules
import React from 'react';
// modules
import { useStore } from 'modules/Stores';
// components
import PagePhotoBlock from 'components/PagePhotoBlock';
import { ArrowForward } from '@mui/icons-material';
import { ContactsAndFormBlock } from 'components/ContactsAndFormBlock';
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
import { PageWrapper } from 'utils/styles';

const CatalogItemScreenView = () => {
  const { getCurrentCatalogItem } = useStore('CatalogItemStore');
  const { setDescription } = useStore('ContactsAndFormBlockStore');
  const searchParams = new URLSearchParams(window.location.search);
  const titleSearchParam = searchParams.get('title');

  console.log('titleSearchParam', titleSearchParam);
  

  const currentItem = getCurrentCatalogItem(titleSearchParam);

  const onOrderButtonClick = (item: PageItems) => {
    setDescription(item.title);

    const anchor = document.querySelector('#contact-form-anchor');

    if (anchor) {
      anchor.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      });
    }
  }
  

  return (
    <PageWrapper>
      { currentItem && (
        <>
          <PagePhotoBlock heading={ currentItem.title } btnText={ currentItem.buttonText } imageUrl={ currentItem.imageUrl } />
          <CatalogWrapper isMainPage>
            <StyledGridContainer container rowSpacing={3} columnSpacing={{ xs: 0, sm: 0, md: 0, lg: 3 }}>
              {currentItem.items.map((item) => (
                <StyledGrid key={ item.id } item xs={12} lg={6}>
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
                        <ItemButton onClick={e => onOrderButtonClick(item)} color="success" variant="text" endIcon={<ArrowForward />}>
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
        </>
      )}
      
      <ContactsAndFormBlock isCatalogItemPage />
    </PageWrapper>
  );
};

export default CatalogItemScreenView;
