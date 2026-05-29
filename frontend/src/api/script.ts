import client from './client';

export async function listScriptFiles(): Promise<string[]> {
  const res: any = await client.get('/script-files');
  return res.data;
}
