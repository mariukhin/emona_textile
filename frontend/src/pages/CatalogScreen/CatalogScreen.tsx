// node modules
import React from 'react';
// components
import { Catalog } from 'components/Catalog';
import PagePhotoBlock from 'components/PagePhotoBlock';
import { PageWrapper } from 'utils/styles';

const CatalogScreenView = () => {
  return (
    <PageWrapper>
      <PagePhotoBlock heading="Каталог" btnText="Якісні товари" imageUrl='assets/hotelno-restor-old.png'/>
      <Catalog />
    </PageWrapper>
  );
};

export default CatalogScreenView;
