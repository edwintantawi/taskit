import React from 'react';

import { Button, Input } from '~/common/components';

export function SignInForm() {
  const handleSubmit = () => {
    console.log('OK');
  };

  return (
    <form className="space-y-3" onSubmit={handleSubmit}>
      <Input required type="email" placeholder="gopher@go.dev" />
      <Input required type="password" placeholder="Secret password" />
      <div className="space-y-4 pt-4">
        <Button fullWidth type="submit" size="large" variants="contained">
          Login account
        </Button>
      </div>
    </form>
  );
}
