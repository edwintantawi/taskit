import React from 'react';
import { Link } from 'react-router-dom';

import { FormLayout } from '~/common/layouts';
import { SignUpForm } from '~/user/containers';

export function SignUpPage() {
  return (
    <FormLayout
      title="Sign-Up Account"
      subtitle="Start a new awsome journey by creating a new account"
      form={<SignUpForm />}
    >
      <p className="text-end text-sm text-gray-500">
        Already have an account?{' '}
        <Link to="/authentications/sign-in" className="text-black underline">
          Sign-In now
        </Link>
      </p>
    </FormLayout>
  );
}
