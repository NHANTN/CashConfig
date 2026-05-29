import { useEffect, useState } from 'react';
import { Table, Card, Select, Tag } from 'antd';
import { listLogs } from '../api/log';
import type { OperationLog as LogType } from '../api/log';

const ACTION_COLORS: Record<string, string> = {
  created: 'green', updated: 'blue', deleted: 'red',
};

export default function OperationLog() {
  const [list, setList] = useState<LogType[]>([]);
  const [loading, setLoading] = useState(false);
  const [action, setAction] = useState<string>();

  const fetch = async () => {
    setLoading(true);
    try {
      const data = await listLogs(action ? { action } : {});
      setList(data);
    } finally { setLoading(false); }
  };

  useEffect(() => { fetch(); }, [action]);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '操作人', dataIndex: 'username', key: 'username', width: 100 },
    {
      title: '操作', dataIndex: 'action', key: 'action', width: 90,
      render: (v: string) => <Tag color={ACTION_COLORS[v]}>{v}</Tag>,
    },
    { title: '目标类型', dataIndex: 'target_type', key: 'target_type', width: 100 },
    { title: '目标 ID', dataIndex: 'target_id', key: 'target_id', width: 80 },
    { title: '详情', dataIndex: 'detail', key: 'detail', ellipsis: true },
    { title: '时间', dataIndex: 'created_at', key: 'created_at', width: 180, render: (v: string) => new Date(v).toLocaleString() },
  ];

  return (
    <Card title="操作日志" extra={
      <Select
        allowClear placeholder="操作类型" value={action} onChange={setAction}
        options={[
          { label: '全部', value: '' },
          { label: '创建', value: 'created' },
          { label: '更新', value: 'updated' },
          { label: '删除', value: 'deleted' },
        ]}
        style={{ width: 120 }}
      />
    }>
      <Table rowKey="id" columns={columns} dataSource={list} loading={loading} pagination={{ pageSize: 20 }} />
    </Card>
  );
}
