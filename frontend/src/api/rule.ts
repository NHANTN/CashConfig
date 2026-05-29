import client from './client';

export interface Rule {
  id: number;
  name: string;
  type: string;
  location: string;
  env_name: string;
  rule: string;
  result: string;
}

export interface RuleListParams {
  type?: string;
  location?: string;
  env_name?: string;
}

export async function listRules(params?: RuleListParams): Promise<Rule[]> {
  const res: any = await client.get('/rules', { params });
  return res.data;
}

export async function getRule(id: number): Promise<Rule> {
  const res: any = await client.get(`/rules/${id}`);
  return res.data;
}

export async function createRule(data: Partial<Rule>): Promise<Rule> {
  const res: any = await client.post('/rules', data);
  return res.data;
}

export async function updateRule(id: number, data: Partial<Rule>): Promise<Rule> {
  const res: any = await client.put(`/rules/${id}`, data);
  return res.data;
}

export async function deleteRule(id: number): Promise<void> {
  await client.delete(`/rules/${id}`);
}

export async function importRuleCsv(file: File): Promise<{ imported: number }> {
  const form = new FormData();
  form.append('file', file);
  const res: any = await client.post('/rules/import', form);
  return res.data;
}

export function getExportRuleCsvUrl(params?: RuleListParams): string {
  const q = new URLSearchParams();
  if (params?.type) q.set('type', params.type);
  if (params?.location) q.set('location', params.location);
  if (params?.env_name) q.set('env_name', params.env_name);
  return `/api/rules/export/csv${q.toString() ? '?' + q.toString() : ''}`;
}
