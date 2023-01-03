import React from 'react';
import { Route, Routes } from 'react-router-dom';

import { AuthGuard, GuestGuard } from '~/authentication/components';
import { SignInPage } from '~/authentication/pages';
import { SignUpPage } from '~/user/pages';
import { HomePage } from '~/common/pages';
import { TaskPage } from '~/task/pages';

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<GuestGuard />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/authentications">
          <Route path="sign-in" element={<SignInPage />} />
          <Route path="sign-up" element={<SignUpPage />} />
        </Route>
      </Route>

      <Route path="/app" element={<AuthGuard />}>
        <Route path="" element={<TaskPage />} />
      </Route>
    </Routes>
  );
}
