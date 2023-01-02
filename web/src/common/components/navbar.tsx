import React from 'react';
import { Link } from 'react-router-dom';

import { Logo, Button } from '~/common/components';

export function Navbar() {
  return (
    <header className="flex items-center justify-between border-b py-3 px-2">
      <Link to="/">
        <Logo />
      </Link>
      <div className="space-x-2">
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
      </div>
    </header>
  );
}
