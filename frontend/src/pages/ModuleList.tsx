import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table, Button, Space, Tag, Popconfirm, message, Select, Card, Upload,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { useModuleStore } from '../store/module';
import { getExportCsvUrl, importModuleCsv, type ModuleListParams } from '../api/module';
import { listAllGroups } from '../api/group';
import type { Group } from '../api/group';

const ENV_OPTIONS = [
  { label: '生产 (pr)', value: 'pr' },
  { label: '预发布 (pp)', value: 'pp' },
  { label: '实验室 (lab)', value: 'lab' },
];

const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];

const envColor: Record<string, string> = { pr: 'red', pp: 'orange', lab: 'blue' };

export default function ModuleList() {
  const navigate = useNavigate();
  const { list, loading, params, page, pageSize, fetch, remove, setPage } = useModuleStore();
  const [env, setEnv] = useState<string>();
  const [location, setLocation] = useState<string>();
  const [groupMap, setGroupMap] = useState<Record<number, string>>({});

  useEffect(() => {
    fetch(params);
    listAllGroups().then((groups: Group[]) => {
      const m: Record<number, string> = {};
      groups.forEach((g) => { m[g.id] = g.name; });
      setGroupMap(m);
    });
  }, []);

  const handleSearch = () => {
    const p: ModuleListParams = {};
    if (env) p.env = env;
    if (location) p.location = location;
    setPage(1, pageSize);
    fetch(p);
  };

  const handleReset = () => { setEnv(undefined); setLocation(undefined); setPage(1, pageSize); fetch({}); };

  const handleDelete = async (id: number) => {
    try { await remove(id); message.success('已删除'); } catch { message.error('删除失败'); }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '名称', dataIndex: 'name', key: 'name', width: 200, ellipsis: true },
    { title: '地区', dataIndex: 'location', key: 'location', width: 80 },
    {
      title: '环境', dataIndex: 'env', key: 'env', width: 100,
      render: (v: string) => <Tag color={envColor[v] || 'default'}>{v}</Tag>,
    },
    {
      title: '引用 Group', dataIndex: 'modules', key: 'modules', width: 300,
      render: (v: string) => {
        let ids: number[] = [];
        try { ids = JSON.parse(v); } catch { return <span>-</span>; }
        return ids.map((id) => (
          <Tag key={id} color="blue" style={{ marginBottom: 2 }}>{groupMap[id] || `#${id}`}</Tag>
        ));
      },
    },
    {
      title: '操作', key: 'action', width: 180,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/modules/${record.id}`)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="Module 配置管理" extra={
      <Space>
        <Select allowClear placeholder="环境" value={env} onChange={setEnv} options={ENV_OPTIONS} style={{ width: 130 }} />
        <Select allowClear placeholder="地区" value={location} onChange={setLocation} options={LOCATIONS.map(l => ({ label: l, value: l }))} style={{ width: 100 }} />
        <Button onClick={handleSearch}>查询</Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    }>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/modules/new')}>
          新增 Module
        </Button>
        <Button icon={<DownloadOutlined />} href={getExportCsvUrl(params)} target="_blank">
          导出 CSV
        </Button>
        <Upload accept=".csv" showUploadList={false} beforeUpload={async (file) => {
          try {
            const result = await importModuleCsv(file);
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
