import { useQuery } from 'react-query';

import { useAuth } from './use-auth';

import {
  AuthenticationService,
  GetAuthUserResponse,
} from '~/authentication/services';

export function useProfile() {
  const { setUser } = useAuth();

  const onSuccess = ({ data }: GetAuthUserResponse) => {
    setUser(data.payload);
  };

  const query = useQuery('auth-profile', AuthenticationService.profile, {
    retry: false,
    refetchOnWindowFocus: false,
    onSuccess,
  });

  return query;
}
