import React from 'react';

import { Button, Input } from '~/common/components';

export function SignUpForm() {
  const handleSubmit = () => {
    console.log('OK');
  };

  return (
    <form className="space-y-3" onSubmit={handleSubmit}>
      <Input required type="text" placeholder="Gopher" />
      <Input required type="email" placeholder="gopher@go.dev" />
      <Input required type="password" placeholder="Secret password" />
      <div className="space-y-4 pt-4">
        <Button fullWidth type="submit" size="large" variants="contained">
          Create Account
        </Button>
      </div>
    </form>
  );
}
