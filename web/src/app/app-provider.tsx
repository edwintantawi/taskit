import React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { BrowserRouter } from 'react-router-dom';

import { AuthContextProvider } from '~/authentication/contexts';

interface AppProviderProps {
  children?: React.ReactNode;
}

const queryClient = new QueryClient();

// Setup all application provider
export function AppProvider({ children }: AppProviderProps) {
  return (
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <AuthContextProvider>{children}</AuthContextProvider>
      </QueryClientProvider>
    </BrowserRouter>
  );
}
