import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';

import { useAuth } from '~/authentication/hooks';

export function GuestGuard() {
  const { user } = useAuth();

  const isAuthenticated = user !== null;
  if (isAuthenticated) {
    return <Navigate to="/app" />;
  }

  return <Outlet />;
}
