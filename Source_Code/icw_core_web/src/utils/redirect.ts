const POST_LOGOUT_REDIRECT_KEY = 'icw_post_logout_redirect';

export function setPostLogoutRedirect(path: string): void {
  sessionStorage.setItem(POST_LOGOUT_REDIRECT_KEY, path);
}

export function getPostLogoutRedirect(): string {
  return sessionStorage.getItem(POST_LOGOUT_REDIRECT_KEY) ?? '';
}

export function clearPostLogoutRedirect(): void {
  sessionStorage.removeItem(POST_LOGOUT_REDIRECT_KEY);
}
