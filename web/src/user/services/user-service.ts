import { AxiosResponse } from 'axios';

import { http, HTTPResponseSuccess } from '~/common/libs';

export type SignUpPayload = { name: string; email: string; password: string };

export type SignUpResponse = AxiosResponse<
  HTTPResponseSuccess<{ id: string; email: string }>
>;

export class UserService {
  static async signUp(payload: SignUpPayload): Promise<SignUpResponse> {
    return http('/users', {
      method: 'POST',
      data: {
        name: payload.name,
        email: payload.email,
        password: payload.password,
      },
    });
  }
}
