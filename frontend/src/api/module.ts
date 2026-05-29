import client from './client';

export interface Module {
  id: number;
  name: string;
  location: string;
  modules: string;
  env: string;
}

export interface ModuleListParams {
  env?: string;
  location?: string;
}

export async function listModules(params?: ModuleListParams): Promise<Module[]> {
  const res: any = await client.get('/modules', { params });
  return res.data;
}

export async function getModule(id: number): Promise<Module> {
  const res: any = await client.get(`/modules/${id}`);
  return res.data;
}

export async function createModule(data: Partial<Module>): Promise<Module> {
  const res: any = await client.post('/modules', data);
  return res.data;
}

export async function updateModule(id: number, data: Partial<Module>): Promise<Module> {
  const res: any = await client.put(`/modules/${id}`, data);
  return res.data;
}

export async function deleteModule(id: number): Promise<void> {
  await client.delete(`/modules/${id}`);
}

export async function importModuleCsv(file: File): Promise<{ imported: number }> {
  const form = new FormData();
  form.append('file', file);
  const res: any = await client.post('/modules/import', form);
  return res.data;
}

export function getExportCsvUrl(params?: ModuleListParams): string {
  const q = new URLSearchParams();
  if (params?.env) q.set('env', params.env);
  if (params?.location) q.set('location', params.location);
  return `/api/modules/export/csv${q.toString() ? '?' + q.toString() : ''}`;
}
