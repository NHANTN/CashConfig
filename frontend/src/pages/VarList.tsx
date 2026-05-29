import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table, Button, Space, Tag, Popconfirm, message, Select, Card, Input, Upload,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, SearchOutlined, UploadOutlined } from '@ant-design/icons';
import { useVarStore } from '../store/var';
import { getExportVarCsvUrl, importVarCsv } from '../api/var';

const ENV_OPTIONS = [
  { label: '生产 (pr)', value: 'pr' },
  { label: '预发布 (pp)', value: 'pp' },
  { label: '实验室 (lab)', value: 'lab' },
];
const envColor: Record<string, string> = { pr: 'red', pp: 'orange', lab: 'blue' };

export default function VarList() {
  const navigate = useNavigate();
  const { list, loading, params, page, pageSize, fetch, remove, setPage } = useVarStore();
  const [env, setEnv] = useState<string>();
  const [varName, setVarName] = useState('');

  useEffect(() => { fetch(params); }, []);

  const handleSearch = () => {
    const p: Record<string, string> = {};
    if (env) p.env = env;
    if (varName) p.var_name = varName;
    setPage(1, pageSize);
    fetch(p);
  };

  const handleReset = () => { setEnv(undefined); setVarName(''); setPage(1, pageSize); fetch({}); };

  const handleDelete = async (id: number) => {
    try { await remove(id); message.success('已删除'); } catch { message.error('删除失败'); }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '变量名', dataIndex: 'var_name', key: 'var_name', width: 260 },
    { title: '值', dataIndex: 'value', key: 'value', ellipsis: true },
    { title: '环境', dataIndex: 'env', key: 'env', width: 80, render: (v: string) => <Tag color={envColor[v]}>{v}</Tag> },
    { title: 'Matcher', dataIndex: 'matcher', key: 'matcher', width: 80, render: (v: string) => <code>{v}</code> },
    {
      title: '操作', key: 'action', width: 120,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/vars/${record.id}`)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="Var 配置管理" extra={
      <Space>
        <Select allowClear placeholder="环境" value={env} onChange={setEnv} options={ENV_OPTIONS} style={{ width: 130 }} />
        <Input placeholder="变量名" value={varName} onChange={e => setVarName(e.target.value)} style={{ width: 180 }} prefix={<SearchOutlined />} />
        <Button onClick={handleSearch}>查询</Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    }>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/vars/new')}>新增 Var</Button>
        <Button icon={<DownloadOutlined />} href={getExportVarCsvUrl(params)} target="_blank">导出 CSV</Button>
        <Upload accept=".csv" showUploadList={false} beforeUpload={async (file) => {
          try {
            const result = await importVarCsv(file);
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