import React, { useState } from 'react';

import { User } from '~/user/entity';
import { AuthContext } from '~/authentication/contexts/auth-context';

interface AuthContextProviderProps {
  children: React.ReactNode;
}

export function AuthContextProvider({ children }: AuthContextProviderProps) {
  const [user, setUser] = useState<User | null>(null);

  return (
    <AuthContext.Provider value={{ user, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}
