// node modules
import React from 'react';
// components
import { Catalog } from 'components/Catalog';
import PagePhotoBlock from 'components/PagePhotoBlock';

const CatalogScreenView = () => {
  return (
    <div>
      <PagePhotoBlock heading="Каталог" btnText="Якісні товари" imageUrl='assets/hotelno-restor-old.png'/>
      <Catalog />
    </div>
  );
};

export default CatalogScreenView;
