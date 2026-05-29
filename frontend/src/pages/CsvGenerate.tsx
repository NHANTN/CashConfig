import { useEffect, useState } from 'react';
import {
  Card, Button, Space, Table, Tag, message, Select, Modal, Input,
} from 'antd';
import { DownloadOutlined, ThunderboltOutlined, DiffOutlined } from '@ant-design/icons';
import {
  generateCsv, downloadCsv, downloadAllCsv,
  getHistory, getDiff,
} from '../api/csv_generate';
import type { GenerationLog } from '../api/csv_generate';

const { TextArea } = Input;

const FILE_TYPES = [
  { label: '全部', value: '' },
  { label: 'Module', value: 'module' },
  { label: 'Rule', value: 'rule' },
  { label: 'Store', value: 'store' },
  { label: 'TillList', value: 'till' },
  { label: 'Var', value: 'var' },
];

export default function CsvGenerate() {
  const [fileType, setFileType] = useState('');
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState<GenerationLog[]>([]);
  const [diffVisible, setDiffVisible] = useState(false);
  const [diffData, setDiffData] = useState<{ from: string; to: string; from_time: string; to_time: string } | null>(null);
  const [diffType, setDiffType] = useState('module');

  const fetchHistory = async () => {
    try {
      const data = await getHistory();
      setLogs(data);
    } catch {
      message.error('获取历史记录失败');
    }
  };

  useEffect(() => { fetchHistory(); }, []);

  const handleGenerate = async () => {
    setLoading(true);
    try {
      const result = await generateCsv(fileType || undefined);
      message.success(`生成成功: ${result.files} 个文件`);
      fetchHistory();
    } catch {
      message.error('生成失败');
    } finally {
      setLoading(false);
    }
  };

  const handleDiff = async () => {
    if (logs.length < 2) {
      message.warning('需要至少 2 条生成记录才能对比');
      return;
    }
    try {
      const result = await getDiff(diffType, logs[logs.length - 1].detail.split(' ')[1], logs[0].detail.split(' ')[1]);
      setDiffData(result);
      setDiffVisible(true);
    } catch {
      message.error('获取差异失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: '生成时间', dataIndex: 'generated_at', key: 'generated_at', width: 180,
      render: (v: string) => new Date(v).toLocaleString(),
    },
    { title: '文件类型', dataIndex: 'file_type', key: 'file_type' },
    { title: '文件数', dataIndex: 'file_count', key: 'file_count', width: 70 },
    { title: '操作人', dataIndex: 'operator', key: 'operator', width: 120 },
    {
      title: '状态', dataIndex: 'status', key: 'status', width: 80,
      render: (v: string) => <Tag color={v === 'success' ? 'green' : 'red'}>{v}</Tag>,
    },
    { title: '详情', dataIndex: 'detail', key: 'detail', ellipsis: true },
  ];

  return (
    <>
      <Card title="CSV 文件生成">
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <Space>
            <Select options={FILE_TYPES} value={fileType} onChange={setFileType} style={{ width: 140 }} />
            <Button type="primary" icon={<ThunderboltOutlined />} loading={loading} onClick={handleGenerate}>
              生成 CSV
            </Button>
            <Button icon={<DownloadOutlined />} onClick={() => downloadAllCsv().catch(() => message.error('下载失败'))}>
              下载全部 (ZIP)
            </Button>
          </Space>
          <Space>
            {FILE_TYPES.filter(t => t.value).map(t => (
              <Button key={t.value} size="small" icon={<DownloadOutlined />} onClick={() => downloadCsv(t.value).catch(() => message.error('下载失败'))}>
                下载 {t.label}
              </Button>
            ))}
          </Space>
        </Space>
      </Card>

      <Card title="生成历史" style={{ marginTop: 16 }} extra={
        <Space>
          <Select options={FILE_TYPES.filter(t => t.value)} value={diffType} onChange={setDiffType} style={{ width: 120 }} />
          <Button icon={<DiffOutlined />} onClick={handleDiff}>对比差异</Button>
        </Space>
      }>
        <Table rowKey="id" columns={columns} dataSource={logs} pagination={{ pageSize: 10 }} />
      </Card>

      <Modal
        title="CSV 差异对比"
        open={diffVisible}
        onCancel={() => setDiffVisible(false)}
        width={900}
        footer={null}
      >
        {diffData && (
          <Space direction="vertical" style={{ width: '100%' }}>
            <div>
              <strong>版本 {diffData.from_time}</strong> vs <strong>版本 {diffData.to_time}</strong>
            </div>
            <div style={{ display: 'flex', gap: 8 }}>
              <TextArea rows={20} value={diffData.from} readOnly style={{ flex: 1, fontFamily: 'monospace', fontSize: 12 }} />
              <TextArea rows={20} value={diffData.to} readOnly style={{ flex: 1, fontFamily: 'monospace', fontSize: 12 }} />
            </div>
          </Space>
        )}
      </Modal>
    </>
  );
}
