import client from './client';

export interface GenerationLog {
  id: number;
  generated_at: string;
  file_type: string;
  file_count: number;
  operator: string;
  status: string;
  detail: string;
}

export async function generateCsv(fileType?: string): Promise<{ timestamp: string; files: number; types: string }> {
  const url = fileType ? `/csv/generate/${fileType}` : '/csv/generate';
  const res: any = await client.post(url);
  return res.data;
}

function downloadBlob(data: Blob, filename: string) {
  const url = URL.createObjectURL(data);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

export async function downloadCsv(fileType: string) {
  const blob = await client.get(`/csv/download/${fileType}`, { responseType: 'blob' }) as unknown as Blob;
  const filename = `${fileType}.csv`;
  downloadBlob(blob, filename);
}

export async function downloadAllCsv() {
  const blob = await client.get('/csv/download/all', { responseType: 'blob' }) as unknown as Blob;
  downloadBlob(blob, 'csv-files.zip');
}

export async function getHistory(fileType?: string): Promise<GenerationLog[]> {
  const url = fileType ? `/csv/history/${fileType}` : '/csv/history';
  const res: any = await client.get(url);
  return res.data;
}

export async function getDiff(fileType: string, from: string, to: string): Promise<{ from: string; to: string; from_time: string; to_time: string }> {
  const res: any = await client.get(`/csv/diff/${fileType}`, { params: { from, to } });
  return res.data;
}
