import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table, Button, Space, Popconfirm, message, Card, Input, Tag, Select, Upload,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, SearchOutlined, UploadOutlined, EyeOutlined } from '@ant-design/icons';
import { useTillListStore } from '../store/till_list';
import { getExportTillListCsvUrl, importTillListCsv } from '../api/till_list';

const ENV_OPTIONS = [
  { label: '生产 (pr)', value: 'pr' },
  { label: '预发布 (pp)', value: 'pp' },
  { label: '实验室 (lab)', value: 'lab' },
];
const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];
const envColor: Record<string, string> = { pr: 'red', pp: 'orange', lab: 'blue' };

export default function TillListPage() {
  const navigate = useNavigate();
  const { list, loading, params, page, pageSize, fetch, remove, setPage } = useTillListStore();
  const [hostName, setHostName] = useState('');
  const [location, setLocation] = useState<string>();
  const [env, setEnv] = useState<string>();

  useEffect(() => { fetch(params); }, []);

  const handleSearch = () => {
    const p: Record<string, string> = {};
    if (hostName) p.host_name = hostName;
    if (location) p.location = location;
    if (env) p.env = env;
    setPage(1, pageSize);
    fetch(p as any);
  };

  const handleReset = () => { setHostName(''); setLocation(undefined); setEnv(undefined); setPage(1, pageSize); fetch({}); };

  const handleDelete = async (id: number) => {
    try { await remove(id); message.success('已删除'); } catch { message.error('删除失败'); }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '主机名', dataIndex: 'host_name', key: 'host_name', width: 130 },
    { title: '地区', dataIndex: 'location', key: 'location', width: 70 },
    { title: '门店', dataIndex: 'store_number', key: 'store_number', width: 80 },
    {
      title: '环境', dataIndex: 'env', key: 'env', width: 70,
      render: (v: string) => v ? <Tag color={envColor[v]}>{v}</Tag> : '-',
    },
    { title: 'IP', dataIndex: 'ip', key: 'ip', width: 130 },
    { title: '硬件型号', dataIndex: 'hardware_model', key: 'hardware_model', width: 120 },
    {
      title: '最近活跃', dataIndex: 'last_seen', key: 'last_seen', width: 160,
      render: (v: string) => v ? new Date(v).toLocaleString() : '-',
    },
    {
      title: '操作', key: 'action', width: 180,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EyeOutlined />} onClick={() => navigate(`/till-lists/${record.id}`)} />
          <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/till-lists/${record.id}/edit`)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="TillList 设备管理" extra={
      <Space>
        <Input placeholder="主机名" value={hostName} onChange={e => setHostName(e.target.value)} style={{ width: 150 }} prefix={<SearchOutlined />} />
        <Input placeholder="IP" style={{ width: 120 }} />
        <Select allowClear placeholder="地区" value={location} onChange={setLocation} options={LOCATIONS.map(l => ({ label: l, value: l }))} style={{ width: 100 }} />
        <Select allowClear placeholder="环境" value={env} onChange={setEnv} options={ENV_OPTIONS} style={{ width: 130 }} />
        <Button onClick={handleSearch}>查询</Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    }>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/till-lists/new')}>新增 TillList</Button>
        <Button icon={<DownloadOutlined />} href={getExportTillListCsvUrl(params as any)} target="_blank">导出 CSV</Button>
        <Upload accept=".csv" showUploadList={false} beforeUpload={async (file) => {
          try {
            const result = await importTillListCsv(file);
            message.success(`导入成功: ${result.imported} 条`);
            fetch(params as any);
          } catch { message.error('导入失败'); }
          return false;
        }}>
          <Button icon={<UploadOutlined />}>导入 CSV</Button>
        </Upload>
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={list}
        loading={loading}
        scroll={{ x: 1200 }}
        pagination={{
          current: page,
          pageSize,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => setPage(p, ps),
        }}
      />
    </Card>
  );
}