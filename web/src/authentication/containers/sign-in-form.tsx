import React from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSearchParams } from 'react-router-dom';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

import { Alert, Button, Input } from '~/common/components';
import { useSignIn } from '~/authentication/hooks';

type SignInFormField = {
  email: string;
  password: string;
};

const validationSchema = yup.object({
  email: yup.string().email().required(),
  password: yup.string().min(6).required(),
});

export function SignInForm() {
  const [searchParams] = useSearchParams();
  const emailFromSignUp = searchParams.get('email') ?? '';

  const { mutate: signIn, error, isError, isLoading } = useSignIn();

  const { register, handleSubmit, formState } = useForm<SignInFormField>({
    resolver: yupResolver(validationSchema),
    defaultValues: { email: emailFromSignUp },
  });

  const onSubmit: SubmitHandler<SignInFormField> = (data) => {
    const { email, password } = data;
    signIn({ email, password });
  };

  return (
    <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
      <Input
        type="email"
        placeholder="gopher@go.dev"
        error={formState.errors.email?.message}
        {...register('email')}
      />
      <Input
        type="password"
        placeholder="Secret password"
        autoFocus={!!emailFromSignUp}
        error={formState.errors.password?.message}
        {...register('password')}
      />

      {isError && (
        <Alert severity="error">
          {error?.response?.data.error ?? error.message}
        </Alert>
      )}

      <div className="space-y-4 pt-4">
        <Button
          fullWidth
          type="submit"
          size="large"
          variants="contained"
          disabled={isLoading}
        >
          {isLoading ? 'Logged in, please wait....' : 'Login account'}
        </Button>
      </div>
    </form>
  );
}
