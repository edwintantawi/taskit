import React from 'react';
import { Link } from 'react-router-dom';

import { SignInForm } from '~/authentication/components/';
import { AuthenticationLayout } from '~/authentication/layouts';

export function SignInPage() {
  return (
    <AuthenticationLayout
      title="Sign-In Account"
      subtitle="Continue your awsome journey by sign-in existing account"
      form={<SignInForm />}
    >
      <p className="text-end text-sm text-gray-500">
        Not have an account?{' '}
        <Link to="/authentications/sign-up" className="text-black underline">
          Create account now
        </Link>
      </p>
    </AuthenticationLayout>
  );
}
