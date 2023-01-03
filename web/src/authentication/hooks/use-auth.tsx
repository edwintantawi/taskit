import { useContext } from 'react';

import { AuthContext } from '~/authentication/contexts';

export function useAuth() {
  const ctx = useContext(AuthContext);

  if (!ctx) {
    throw new Error('Auth context value is not provided!');
  }

  return {
    user: ctx.user,
    setUser: ctx.setUser,
    isAuthenticate: ctx.user !== null,
  };
}
