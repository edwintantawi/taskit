import React from 'react';
import { Link } from 'react-router-dom';

import { useAuth } from '~/authentication/hooks';
import { useSignOut } from '~/authentication/hooks/use-sign-out';
import { Logo, Button } from '~/common/components';

export function Navbar() {
  const { isAuthenticate } = useAuth();
  const { mutate: signOut, isLoading } = useSignOut();

  const handleSignOut = () => signOut();

  return (
    <header className="flex items-center justify-between border-b py-3 px-2">
      <Link to="/">
        <Logo />
      </Link>
      <div className="space-x-2">
        {!isAuthenticate ? (
          <>
            <Button
              as={Link}
              size="small"
              variants="outlined"
              to="/authentications/sign-in"
            >
              Sign-In
            </Button>
            <Button
              as={Link}
              size="small"
              variants="contained"
              to="/authentications/sign-up"
            >
              Create Account
            </Button>
          </>
        ) : (
          <Button
            size="small"
            variants="contained"
            disabled={isLoading}
            onClick={handleSignOut}
          >
            {isLoading ? 'Logged out, Please wait...' : 'Sign-Out Account'}
          </Button>
        )}
      </div>
    </header>
  );
}
