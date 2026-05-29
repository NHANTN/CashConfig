import client from './client';

export interface Store {
  id: number;
  store_number: string;
  network_segment: string;
  webpos_env: string;
  eft: string;
  location: string;
  rf_server: string;
  cashtill_seg_gw: string;
}

export interface StoreListParams {
  location?: string;
  eft?: string;
}

export async function listStores(params?: StoreListParams): Promise<Store[]> {
  const res: any = await client.get('/stores', { params });
  return res.data;
}

export async function getStore(id: number): Promise<Store> {
  const res: any = await client.get(`/stores/${id}`);
  return res.data;
}

export async function createStore(data: Partial<Store>): Promise<Store> {
  const res: any = await client.post('/stores', data);
  return res.data;
}

export async function updateStore(id: number, data: Partial<Store>): Promise<Store> {
  const res: any = await client.put(`/stores/${id}`, data);
  return res.data;
}

export async function deleteStore(id: number): Promise<void> {
  await client.delete(`/stores/${id}`);
}

export async function importStoreCsv(file: File): Promise<{ imported: number }> {
  const form = new FormData();
  form.append('file', file);
  const res: any = await client.post('/stores/import', form);
  return res.data;
}

export function getExportStoreCsvUrl(params?: StoreListParams): string {
  const q = new URLSearchParams();
  if (params?.location) q.set('location', params.location);
  if (params?.eft) q.set('eft', params.eft);
  return `/api/stores/export/csv${q.toString() ? '?' + q.toString() : ''}`;
}
