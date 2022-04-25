// node modules
import React, { Suspense as ReactSuspense } from 'react';

const LazyLoader = () => (
  <div>
    {/* <Loading /> */}
  </div>
);

const Suspense: React.FC = ({ children }) => (
  <ReactSuspense
    fallback={<LazyLoader />}
  >
    {children}
  </ReactSuspense>
);

export default Suspense;
