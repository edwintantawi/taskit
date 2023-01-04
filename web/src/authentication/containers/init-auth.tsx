import React from 'react';

import { useProfile } from '~/authentication/hooks';

interface InitAuthProps {
  children?: React.ReactNode;
}

export function InitAuth({ children }: InitAuthProps) {
  const { isLoading } = useProfile();

  if (isLoading) {
    return <p>Loading....</p>;
  }
  return <>{children}</>;
}
