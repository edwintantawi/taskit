import { useQuery } from 'react-query';

import { useAuth } from './use-auth';

import {
  AuthenticationService,
  GetAuthUserResponse,
} from '~/authentication/services';

export function useAuthProfile() {
  const { setUser } = useAuth();

  const onSuccess = ({ data }: GetAuthUserResponse) => {
    setUser(data.payload);
  };

  const query = useQuery('auth-profile', AuthenticationService.getAuthProfile, {
    retry: false,
    onSuccess,
  });

  return query;
}
