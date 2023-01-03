import Axios from 'axios';

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
  baseURL: 'http://localhost:5000/api',
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
