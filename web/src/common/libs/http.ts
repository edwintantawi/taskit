import Axios from 'axios';

import { config } from '~/config';
import { TokenManager } from '~/common/services';
import { AuthenticationService } from '~/authentication/services';

export type HTTPResponseSuccess<T> = {
  status_code: number;
  message: string;
  payload: T;
};

export type HTTPResponseError = {
  status_code: number;
  message: string;
  error: string;
};

const http = Axios.create({
  baseURL: config.BASE_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

http.interceptors.request.use(
  (config) => {
    const accessToken = TokenManager.getAccessToken();
    if (accessToken) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      config.headers['Authorization'] = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

http.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = TokenManager.getRefreshToken() ?? '';
        const result = await AuthenticationService.refresh({ refreshToken });

        TokenManager.storeToken({
          access_token: result.data.payload.access_token,
          refresh_token: result.data.payload.refresh_token,
        });

        originalRequest.headers[
          'Authorization'
        ] = `Bearer ${result.data.payload.access_token}`;

        return http(originalRequest);
      } catch (_error) {
        return Promise.reject(_error);
      }
    } else {
      return Promise.reject(error.response.data);
    }
  }
);

export { http };
