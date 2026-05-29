import client from './client';

export interface User {
  id: number;
  username: string;
  name: string;
  role_id: number;
  status: number;
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  const res: any = await client.post('/auth/login', { username, password });
  return res.data;
}

export async function loginLDAP(username: string, password: string): Promise<LoginResponse> {
  const res: any = await client.post('/auth/login/ldap', { username, password });
  return res.data;
}

export async function getSSOLoginURL(): Promise<string> {
  const res: any = await client.get('/auth/sso/login');
  return res.data.auth_url;
}

export async function handleSSOCallback(code: string, state: string): Promise<LoginResponse> {
  const res: any = await client.get(`/auth/sso/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`);
  return res.data;
}

export async function getPermissions(): Promise<string[]> {
  const res: any = await client.get('/auth/permissions');
  return res.data;
}

export async function refreshToken(): Promise<string> {
  const res: any = await client.post('/auth/refresh');
  return res.data.token;
}
