import React from 'react';
import { Link } from 'react-router-dom';

import { Logo, Button } from '~/common/components';

export function Navbar() {
  return (
    <header className="flex items-center justify-between border-b py-3 px-2">
      <Logo />
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
          Sign-Up
        </Button>
      </div>
    </header>
  );
}
