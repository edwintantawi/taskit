import Axios from 'axios';

import { appConfig } from '~/app';
import { TokenManager } from '~/common/services';

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
  baseURL: appConfig.BASE_API_URL,
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

export { http };
