import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClientProvider } from 'react-query';
import { ReactQueryDevtools } from 'react-query/devtools';

import { AuthContextProvider } from '~/authentication/contexts';
import { InitAuthProfile } from '~/authentication/containers/init-auth-profile';
import { queryClient } from '~/common/libs';

interface AppProviderProps {
  children?: React.ReactNode;
}

// Setup all application provider
export function AppProvider({ children }: AppProviderProps) {
  return (
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <AuthContextProvider>
          <InitAuthProfile>{children}</InitAuthProfile>
        </AuthContextProvider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </BrowserRouter>
  );
}
