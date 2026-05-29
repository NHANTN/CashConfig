import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table, Button, Space, Popconfirm, message, Card, Input, Tag,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import { useGroupStore } from '../store/group';

export default function GroupList() {
  const navigate = useNavigate();
  const { list, loading, page, pageSize, fetch, remove, setPage } = useGroupStore();
  const [name, setName] = useState('');

  useEffect(() => { fetch(); }, []);

  const handleSearch = () => { setPage(1, pageSize); fetch(name || undefined); };
  const handleReset = () => { setName(''); setPage(1, pageSize); fetch(); };

  const handleDelete = async (id: number) => {
    try { await remove(id); message.success('已删除'); }
    catch { message.error('删除失败'); }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '组名称', dataIndex: 'name', key: 'name' },
    {
      title: 'Steps', dataIndex: 'steps', key: 'steps', ellipsis: true,
      render: (v: string) => {
        try {
          const steps = JSON.parse(v);
          return steps.map((s: any, i: number) => (
            <Tag key={i} style={{ marginBottom: 2 }}>{s.name}</Tag>
          ));
        } catch { return <code>{v?.substring(0, 60)}</code>; }
      },
    },
    {
      title: '操作', key: 'action', width: 120,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/groups/${record.id}`)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="Group 组管理" extra={
      <Space>
        <Input placeholder="组名称" value={name} onChange={e => setName(e.target.value)} style={{ width: 180 }} prefix={<SearchOutlined />} />
        <Button onClick={handleSearch}>查询</Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    }>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/groups/new')}>新增 Group</Button>
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