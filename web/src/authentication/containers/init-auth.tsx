import React from 'react';

import { useProfile } from '~/authentication/hooks';
import { LoadingScreen } from '~/common/components';

interface InitAuthProps {
  children?: React.ReactNode;
}

export function InitAuth({ children }: InitAuthProps) {
  const { isLoading } = useProfile();

  if (isLoading) {
    return <LoadingScreen />;
  }
  return <>{children}</>;
}
