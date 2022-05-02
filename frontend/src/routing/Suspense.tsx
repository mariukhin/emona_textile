// node modules
import React, { Suspense as ReactSuspense } from 'react';
import { observer } from 'mobx-react';

interface SuspenseProps {
  children: React.ReactNode;
}

const LazyLoader = () => (
  <div>
    {/* <Loading /> */}
  </div>
);

const Suspense: React.FC<SuspenseProps> = ({ children }) => (
  <ReactSuspense
    fallback={<LazyLoader />}
  >
    {children}
  </ReactSuspense>
);

export default observer(Suspense);
