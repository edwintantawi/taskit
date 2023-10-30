import React from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { SubmitHandler, useForm } from 'react-hook-form';
import * as yup from 'yup';

import { Alert, Button, Input } from '~/common/components';
import { useSignUp } from '~/user/hooks';

type SignUpFormField = {
  name: string;
  email: string;
  password: string;
};

const validationSchema = yup.object({
  name: yup.string().required(),
  email: yup.string().email().required(),
  password: yup.string().min(6).required(),
});

export function SignUpForm() {
  const { mutate: signUp, error, isError, isLoading } = useSignUp();
  const { register, handleSubmit, formState } = useForm<SignUpFormField>({
    resolver: yupResolver(validationSchema),
  });

  const onSubmit: SubmitHandler<SignUpFormField> = (data) => {
    const { name, email, password } = data;
    signUp({ name, email, password });
  };

  return (
    <form noValidate className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
      <Input
        type="text"
        placeholder="Gopher"
        error={formState.errors.name?.message}
        {...register('name')}
      />
      <Input
        type="email"
        placeholder="gopher@go.dev"
        error={formState.errors.email?.message}
        {...register('email')}
      />
      <Input
        type="password"
        placeholder="Secret password"
        error={formState.errors.password?.message}
        {...register('password')}
      />

      {isError && <Alert severity="error">{error.error}</Alert>}

      <div className="space-y-4 pt-4">
        <Button
          fullWidth
          type="submit"
          size="large"
          variants="contained"
          disabled={isLoading}
        >
          {isLoading ? 'Creating, please wait....' : 'Create account'}
        </Button>
      </div>
    </form>
  );
}
