import { useMutation } from 'react-query';
import { useNavigate } from 'react-router-dom';

import {
  AuthenticationService,
  SignInPayload,
  SignInResponse,
} from '~/authentication/services';
import { TokenManager } from '~/common/services';
import { HTTPResponseError, queryClient } from '~/common/libs';

export function useSignIn() {
  const navigate = useNavigate();

  const onSuccess = ({ data }: SignInResponse) => {
    TokenManager.storeToken({
      access_token: data.payload.access_token,
      refresh_token: data.payload.refresh_token,
    });
    queryClient.invalidateQueries({ queryKey: ['auth-profile'] });
    navigate('/app');
  };

  const mutation = useMutation<
    SignInResponse,
    HTTPResponseError,
    SignInPayload
  >(AuthenticationService.signIn, { onSuccess });

  return mutation;
}
