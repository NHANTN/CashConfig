import client from './client';

export interface OperationLog {
  id: number;
  user_id: number;
  username: string;
  action: string;
  target_type: string;
  target_id: string;
  detail: string;
  created_at: string;
}

export async function listLogs(params?: { target_type?: string; action?: string }): Promise<OperationLog[]> {
  const res: any = await client.get('/system/logs', { params });
  return res.data;
}
