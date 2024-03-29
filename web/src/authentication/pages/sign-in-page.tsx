import React from 'react';
import { Link } from 'react-router-dom';

import { SignInForm } from '~/authentication/containers';
import { FormLayout } from '~/common/layouts';

export function SignInPage() {
  return (
    <FormLayout
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
    </FormLayout>
  );
}
