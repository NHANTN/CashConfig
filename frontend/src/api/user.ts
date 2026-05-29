import client from './client';

export interface User {
  id: number;
  username: string;
  name: string;
  role_id: number;
  status: number;
  created_at: string;
  updated_at: string;
  role?: { id: number; name: string; code: string };
}

export async function listUsers(): Promise<User[]> {
  const res: any = await client.get('/system/users');
  return res.data;
}

export async function getUser(id: number): Promise<User> {
  const res: any = await client.get(`/system/users/${id}`);
  return res.data;
}

export async function createUser(data: { username: string; password: string; name: string; role_id: number }): Promise<User> {
  const res: any = await client.post('/system/users', data);
  return res.data;
}

export async function updateUser(id: number, data: { username?: string; name?: string; role_id?: number; status?: number; password?: string }): Promise<void> {
  await client.put(`/system/users/${id}`, data);
}

export async function deleteUser(id: number): Promise<void> {
  await client.delete(`/system/users/${id}`);
}
