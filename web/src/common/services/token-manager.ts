import { Authentication } from '~/authentication/entity';

const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

export class TokenManager {
  static storeToken({ access_token, refresh_token }: Authentication) {
    localStorage.setItem(ACCESS_TOKEN_KEY, access_token);
    localStorage.setItem(REFRESH_TOKEN_KEY, refresh_token);
  }
  static getAccessToken() {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }
  static getRefreshToken() {
    return localStorage.getItem(REFRESH_TOKEN_KEY);
  }
}
