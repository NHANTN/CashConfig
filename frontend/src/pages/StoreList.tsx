import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table, Button, Space, Popconfirm, message, Select, Card, Upload,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { useStoreStore } from '../store/store';
import { getExportStoreCsvUrl, importStoreCsv } from '../api/store';

const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];

export default function StoreList() {
  const navigate = useNavigate();
  const { list, loading, params, page, pageSize, fetch, remove, setPage } = useStoreStore();
  const [location, setLocation] = useState<string>();

  useEffect(() => { fetch(params); }, []);

  const handleSearch = () => {
    const p: Record<string, string> = {};
    if (location) p.location = location;
    setPage(1, pageSize);
    fetch(p);
  };

  const handleReset = () => { setLocation(undefined); setPage(1, pageSize); fetch({}); };

  const handleDelete = async (id: number) => {
    try { await remove(id); message.success('已删除'); } catch { message.error('删除失败'); }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '门店号', dataIndex: 'store_number', key: 'store_number', width: 80 },
    { title: '网段', dataIndex: 'network_segment', key: 'network_segment', width: 130 },
    { title: '环境', dataIndex: 'webpos_env', key: 'webpos_env', width: 70 },
    { title: 'EFT', dataIndex: 'eft', key: 'eft', width: 100, ellipsis: true },
    { title: '地区', dataIndex: 'location', key: 'location', width: 70 },
    { title: 'RF 服务器', dataIndex: 'rf_server', key: 'rf_server', width: 140 },
    { title: 'Cashtill 网关', dataIndex: 'cashtill_seg_gw', key: 'cashtill_seg_gw', width: 140 },
    {
      title: '操作', key: 'action', width: 120,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/stores/${record.id}`)} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="Store 配置管理" extra={
      <Space>
        <Select allowClear placeholder="地区" value={location} onChange={setLocation} options={LOCATIONS.map(l => ({ label: l, value: l }))} style={{ width: 100 }} />
        <Button onClick={handleSearch}>查询</Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    }>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/stores/new')}>新增 Store</Button>
        <Button icon={<DownloadOutlined />} href={getExportStoreCsvUrl(params)} target="_blank">导出 CSV</Button>
        <Upload accept=".csv" showUploadList={false} beforeUpload={async (file) => {
          try {
            const result = await importStoreCsv(file);
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
        scroll={{ x: 1000 }}
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