import React from 'react';
import { Route, Routes } from 'react-router-dom';

import { SignInPage, SignUpPage } from '~/authentication/pages';
import { HomePage } from '~/common/pages';

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/authentications">
        <Route path="sign-in" element={<SignInPage />} />
        <Route path="sign-up" element={<SignUpPage />} />
      </Route>
    </Routes>
  );
}
