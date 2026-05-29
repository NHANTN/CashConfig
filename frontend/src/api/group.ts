import client from './client';

export interface Group {
  id: number;
  name: string;
  steps: string;
}

export async function listGroups(name?: string): Promise<Group[]> {
  const params = name ? { name } : undefined;
  const res: any = await client.get('/groups', { params });
  return res.data;
}

export async function listAllGroups(): Promise<Group[]> {
  const res: any = await client.get('/groups/all');
  return res.data;
}

export async function getGroup(id: number): Promise<Group> {
  const res: any = await client.get(`/groups/${id}`);
  return res.data;
}

export async function createGroup(data: Partial<Group>): Promise<Group> {
  const res: any = await client.post('/groups', data);
  return res.data;
}

export async function updateGroup(id: number, data: Partial<Group>): Promise<Group> {
  const res: any = await client.put(`/groups/${id}`, data);
  return res.data;
}

export async function deleteGroup(id: number): Promise<void> {
  await client.delete(`/groups/${id}`);
}
