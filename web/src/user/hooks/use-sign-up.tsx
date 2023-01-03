import { AxiosError } from 'axios';
import { useMutation } from 'react-query';
import { useNavigate } from 'react-router-dom';

import { UserService, SignUpPayload, SignUpResponse } from '~/user/services';
import { HTTPResponseError } from '~/common/libs';

export function useSignUp() {
  const navigate = useNavigate();

  const onSuccess = ({ data }: SignUpResponse) => {
    const { email } = data.payload;
    navigate(`/authentications/sign-in?email=${email}`);
  };

  const mutation = useMutation<
    SignUpResponse,
    AxiosError<HTTPResponseError>,
    SignUpPayload
  >(UserService.signUp, { onSuccess });

  return mutation;
}
