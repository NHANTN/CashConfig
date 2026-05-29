import client from './client';

export interface Var {
  id: number;
  var_name: string;
  value: string;
  env: string;
  matcher: string;
}

export interface VarListParams {
  env?: string;
  var_name?: string;
}

export async function listVars(params?: VarListParams): Promise<Var[]> {
  const res: any = await client.get('/vars', { params });
  return res.data;
}

export async function getVar(id: number): Promise<Var> {
  const res: any = await client.get(`/vars/${id}`);
  return res.data;
}

export async function createVar(data: Partial<Var>): Promise<Var> {
  const res: any = await client.post('/vars', data);
  return res.data;
}

export async function updateVar(id: number, data: Partial<Var>): Promise<Var> {
  const res: any = await client.put(`/vars/${id}`, data);
  return res.data;
}

export async function deleteVar(id: number): Promise<void> {
  await client.delete(`/vars/${id}`);
}

export async function importVarCsv(file: File): Promise<{ imported: number }> {
  const form = new FormData();
  form.append('file', file);
  const res: any = await client.post('/vars/import', form);
  return res.data;
}

export function getExportVarCsvUrl(params?: VarListParams): string {
  const q = new URLSearchParams();
  if (params?.env) q.set('env', params.env);
  if (params?.var_name) q.set('var_name', params.var_name);
  return `/api/vars/export/csv${q.toString() ? '?' + q.toString() : ''}`;
}
