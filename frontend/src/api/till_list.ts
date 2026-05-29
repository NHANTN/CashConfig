import client from './client';

export interface TillList {
  id: number;
  host_name: string;
  mac_address: string;
  location: string;
  store_number: string;
  env: string;
  name: string;
  ip: string;
  hardware_model: string;
  group_id: number;
  request_body: string;
  last_seen: string;
}

export interface TillListParams {
  host_name?: string;
  location?: string;
  env?: string;
}

export async function listTillLists(params?: TillListParams): Promise<TillList[]> {
  const res: any = await client.get('/till-lists', { params });
  return res.data;
}

export async function getTillList(id: number): Promise<TillList> {
  const res: any = await client.get(`/till-lists/${id}`);
  return res.data;
}

export async function createTillList(data: Partial<TillList>): Promise<TillList> {
  const res: any = await client.post('/till-lists', data);
  return res.data;
}

export async function updateTillList(id: number, data: Partial<TillList>): Promise<TillList> {
  const res: any = await client.put(`/till-lists/${id}`, data);
  return res.data;
}

export async function deleteTillList(id: number): Promise<void> {
  await client.delete(`/till-lists/${id}`);
}

export async function importTillListCsv(file: File): Promise<{ imported: number }> {
  const form = new FormData();
  form.append('file', file);
  const res: any = await client.post('/till-lists/import', form);
  return res.data;
}

export function getExportTillListCsvUrl(params?: TillListParams): string {
  const q = new URLSearchParams();
  if (params?.host_name) q.set('host_name', params.host_name);
  if (params?.location) q.set('location', params.location);
  if (params?.env) q.set('env', params.env);
  return `/api/till-lists/export/csv${q.toString() ? '?' + q.toString() : ''}`;
}

export interface SyncReport {
  id: number;
  till_list_id: number;
  request_body: string;
  module_execution: string;
  status: number;
  duration: number;
  sync_time: string;
  created_at: string;
}

export interface ModuleStep {
  Name: string;
  Status: boolean;
  Duration: number;
  Output: string;
}

export interface ModuleExec {
  Name: string;
  Status: boolean;
  Duration: number;
  Steps: ModuleStep[];
}

export async function listSyncReports(tillListId: number): Promise<SyncReport[]> {
  const res: any = await client.get(`/till-lists/${tillListId}/reports`);
  return res.data;
}

export interface DeviceWithReports {
  till_list: TillList;
  reports: SyncReport[];
}

export async function queryReportsByDevice(params: { host_name?: string; mac_address?: string }): Promise<DeviceWithReports[]> {
  const res: any = await client.get('/till-lists/reports', { params });
  return res.data;
}

export async function getSyncReport(tillListId: number, reportId: number): Promise<SyncReport> {
  const res: any = await client.get(`/till-lists/${tillListId}/reports/${reportId}`);
  return res.data;
}

export interface RequestBodyJson {
  Ip: string;
  TillNumber: string;
  Execution: {
    EndTime: string;
    Status: number;
    StartTime: string;
    ModuleExecution: Array<{
      Duration: number;
      Name: string;
      Status: boolean;
      Steps: Array<{
        Output: string;
        Status: boolean;
        Name: string;
        Duration: number;
      }>;
    }>;
    Duration: number;
  };
  CreationTime: string;
  Location: string;
  Group: number;
  Fact: Array<{
    Key: string;
    Value: string;
    Type: number;
  }>;
  StoreNumber: string;
  Name: string;
  Env: string;
  HardwareModel: string;
}