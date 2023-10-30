import React from 'react';

import { AppRouter } from '~/app';
import { Navbar } from '~/common/components';
import { AppLayout, MainLayout } from '~/common/layouts';

export function App() {
  return (
    <AppLayout>
      <Navbar />
      <MainLayout>
        <AppRouter />
      </MainLayout>
    </AppLayout>
  );
}
