import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClientProvider } from 'react-query';
import { ReactQueryDevtools } from 'react-query/devtools';

import { AuthContextProvider } from '~/authentication/contexts';
import { InitAuth } from '~/authentication/containers';
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
          <InitAuth>{children}</InitAuth>
        </AuthContextProvider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </BrowserRouter>
  );
}
