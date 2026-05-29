import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table, Button, Space, Tag, Popconfirm, message, Select, Card, Upload,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { useRuleStore } from '../store/rule';
import { getExportRuleCsvUrl, importRuleCsv } from '../api/rule';

const TYPES = [
  { label: '环境 (env)', value: 'env' },
  { label: '分组 (group)', value: 'group' },
];
const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];

const typeColor: Record<string, string> = { env: 'blue', group: 'green' };

export default function RuleList() {
  const navigate = useNavigate();
  const { list, loading, params, page, pageSize, fetch, remove, setPage } = useRuleStore();
  const [type, setType] = useState<string>();
  const [location, setLocation] = useState<string>();

  useEffect(() => { fetch(params); }, []);

  const handleSearch = () => {
    const p: Record<string, string> = {};
    if (type) p.type = type;
    if (location) p.location = location;
    setPage(1, pageSize);
    fetch(p);
  };

  const handleReset = () => { setType(undefined); setLocation(undefined); setPage(1, pageSize); fetch({}); };

  const handleDelete = async (id: number) => {
    try { await remove(id); message.success('已删除'); } catch { message.error('删除失败'); }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '名称', dataIndex: 'name', key: 'name', ellipsis: true },
    { title: '类型', dataIndex: 'type', key: 'type', width: 80, render: (v: string) => <Tag color={typeColor[v]}>{v}</Tag> },
    { title: '地区', dataIndex: 'location', key: 'location', width: 70 },
    { title: '环境变量', dataIndex: 'env_name', key: 'env_name', width: 180, ellipsis: true },
    { title: '规则', dataIndex: 'rule', key: 'rule', ellipsis: true },
    { title: '结果', dataIndex: 'result', key: 'result', ellipsis: true },
    {
      title: '操作', key: 'action', width: 120,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/rules/${record.id}`)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="Rule 配置管理" extra={
      <Space>
        <Select allowClear placeholder="类型" value={type} onChange={setType} options={TYPES} style={{ width: 120 }} />
        <Select allowClear placeholder="地区" value={location} onChange={setLocation} options={LOCATIONS.map(l => ({ label: l, value: l }))} style={{ width: 100 }} />
        <Button onClick={handleSearch}>查询</Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    }>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/rules/new')}>新增 Rule</Button>
        <Button icon={<DownloadOutlined />} href={getExportRuleCsvUrl(params)} target="_blank">导出 CSV</Button>
        <Upload accept=".csv" showUploadList={false} beforeUpload={async (file) => {
          try {
            const result = await importRuleCsv(file);
            message.success(`导入成功: ${result.imported} 条`);
            fetch(params);
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