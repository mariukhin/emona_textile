// node modules
import React from 'react';
// components
import Catalog from 'components/Catalog';
import PagePhotoBlock from 'components/PagePhotoBlock';
// mocks
import { mockedCatalogItems } from 'components/Catalog/mocks';

const CatalogScreenView = () => {
  return (
    <div>
      <PagePhotoBlock heading="Каталог" btnText="Якісні товари" imageUrl='assets/hotelno-restor.jpeg'/>
      <Catalog catalogItems={mockedCatalogItems} />
    </div>
  );
};

export default CatalogScreenView;
