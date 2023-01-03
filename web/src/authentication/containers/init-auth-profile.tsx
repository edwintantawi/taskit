import React from 'react';

import { useAuthProfile } from '~/authentication/hooks';

interface InitAuthProfileProps {
  children?: React.ReactNode;
}

export function InitAuthProfile({ children }: InitAuthProfileProps) {
  const { isLoading } = useAuthProfile();

  if (isLoading) {
    return <p>Loading....</p>;
  }
  return <>{children}</>;
}
