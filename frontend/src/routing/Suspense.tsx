// node modules
import React, { Suspense as ReactSuspense } from 'react';
import { observer } from 'mobx-react';
import { TailSpin } from  'react-loader-spinner';
import { colors } from 'utils/color';

interface SuspenseProps {
  children: React.ReactNode;
}

const LazyLoader = () => (
  <div style={{ width: '100%', height: '800px', display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
    <TailSpin
      height="100"
      width="100"
      color={colors.background.green}
      ariaLabel='loading'
    />
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
