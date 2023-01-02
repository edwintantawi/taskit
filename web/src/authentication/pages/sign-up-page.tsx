import React from 'react';
import { Link } from 'react-router-dom';

import { AuthenticationLayout } from '~/authentication/layouts';
import { SignUpForm } from '~/authentication/components';

export function SignUpPage() {
  return (
    <AuthenticationLayout
      title="Sign-Up Account"
      subtitle="Start a new journey by creating a new account"
      form={<SignUpForm />}
    >
      <p className="text-end text-sm text-gray-500">
        Already have an account?{' '}
        <Link to="/authentications/sign-in" className="text-black underline">
          Sign-In now
        </Link>
      </p>
    </AuthenticationLayout>
  );
}
