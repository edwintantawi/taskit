import React from 'react';
import { Route, Routes } from 'react-router-dom';

import { SignInPage, SignUpPage } from '../authentication/pages';

function AppRouter() {
  return (
    <Routes>
      <Route path="/authentications">
        <Route path="sign-in" element={<SignInPage />} />
        <Route path="sign-up" element={<SignUpPage />} />
      </Route>
    </Routes>
  );
}

export { AppRouter };
