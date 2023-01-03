import { AxiosResponse } from 'axios';

import { http, HTTPResponseSuccess } from '~/common/libs';
import { Authentication } from '~/authentication/entity';

export type SignInPayload = { email: string; password: string };
export type SignOutPayload = { refreshToken: string };
export type RefreshPayload = { refreshToken: string };

export type SignInResponse = AxiosResponse<HTTPResponseSuccess<Authentication>>;
export type RefreshResponse = AxiosResponse<
  HTTPResponseSuccess<Authentication>
>;

export class AuthenticationService {
  static async signIn(payload: SignInPayload): Promise<SignInResponse> {
    return http('/authentications', {
      method: 'POST',
      data: {
        email: payload.email,
        password: payload.password,
      },
    });
  }

  static async refresh(payload: RefreshPayload): Promise<RefreshResponse> {
    return http('/authentications', {
      method: 'PUT',
      data: {
        refresh_token: payload.refreshToken,
      },
    });
  }

  static async signOut(payload: SignOutPayload) {
    return http('/authentications', {
      method: 'DELETE',
      data: {
        refresh_token: payload.refreshToken,
      },
    });
  }
}
