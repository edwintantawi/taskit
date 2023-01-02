import React from 'react';

import { AppRouter } from '~/app';
import { Navbar } from '~/common/components';
import { AppLayout } from '~/common/layouts';
import { MainLayout } from '~/common/layouts/main-layout';

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
