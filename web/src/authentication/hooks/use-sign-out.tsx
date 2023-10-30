import { useMutation } from 'react-query';
import { useNavigate } from 'react-router-dom';

import {
  AuthenticationService,
  SignOutPayload,
} from '~/authentication/services';
import { TokenManager } from '~/common/services';
import { HTTPResponseError } from '~/common/libs';
import { useAuth } from '~/authentication/hooks';

export function useSignOut() {
  const navigate = useNavigate();
  const { setUser } = useAuth();

  const onSuccess = () => {
    TokenManager.clearToken();
    setUser(null);
    navigate('/');
  };

  const mutation = useMutation<unknown, HTTPResponseError, SignOutPayload>(
    AuthenticationService.signOut,
    { onSuccess, onError: (error) => alert(error.error) }
  );

  const signOut = () => {
    mutation.mutate({
      refreshToken: TokenManager.getRefreshToken() ?? '',
    });
  };

  return { ...mutation, mutate: signOut };
}
