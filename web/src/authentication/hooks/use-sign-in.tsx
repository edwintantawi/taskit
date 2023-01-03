import { AxiosError } from 'axios';
import { useMutation } from 'react-query';
import { useNavigate } from 'react-router-dom';

import {
  AuthenticationService,
  SignInPayload,
  SignInResponse,
} from '~/authentication/services';
import { TokenManager } from '~/common/services';
import { HTTPResponseError } from '~/common/libs';

export function useSignIn() {
  const navigate = useNavigate();

  const onSuccess = ({ data }: SignInResponse) => {
    TokenManager.storeToken({
      access_token: data.payload.access_token,
      refresh_token: data.payload.refresh_token,
    });
    navigate('/app');
  };

  const mutation = useMutation<
    SignInResponse,
    AxiosError<HTTPResponseError>,
    SignInPayload
  >(AuthenticationService.signIn, { onSuccess });

  return mutation;
}
