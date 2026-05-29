import client from './client';

export interface Role {
  id: number;
  name: string;
  code: string;
  permissions: string;
  status: number;
  created_at: string;
}

export async function listRoles(): Promise<Role[]> {
  const res: any = await client.get('/system/roles');
  return res.data;
}

export async function getRole(id: number): Promise<Role> {
  const res: any = await client.get(`/system/roles/${id}`);
  return res.data;
}

export async function createRole(data: { name: string; code: string; permissions: string }): Promise<Role> {
  const res: any = await client.post('/system/roles', data);
  return res.data;
}

export async function updateRole(id: number, data: { name?: string; code?: string; permissions?: string; status?: number }): Promise<Role> {
  const res: any = await client.put(`/system/roles/${id}`, data);
  return res.data;
}

export async function deleteRole(id: number): Promise<void> {
  await client.delete(`/system/roles/${id}`);
}
