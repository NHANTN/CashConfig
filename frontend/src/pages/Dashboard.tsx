import { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic, Spin } from 'antd';
import {
  AppstoreOutlined, NodeIndexOutlined, ShopOutlined,
  LaptopOutlined, CodeOutlined, FolderOutlined, UserOutlined,
} from '@ant-design/icons';
import client from '../api/client';

interface Stats {
  module_count: number;
  rule_count: number;
  store_count: number;
  till_count: number;
  var_count: number;
  group_count: number;
  user_count: number;
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    client.get('/dashboard/stats').then((res: any) => {
      setStats(res.data);
    }).finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />;

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col span={6}><Card><Statistic title="Module 配置" value={stats?.module_count} prefix={<AppstoreOutlined />} /></Card></Col>
        <Col span={6}><Card><Statistic title="Rule 规则" value={stats?.rule_count} prefix={<NodeIndexOutlined />} /></Card></Col>
        <Col span={6}><Card><Statistic title="Store 门店" value={stats?.store_count} prefix={<ShopOutlined />} /></Card></Col>
        <Col span={6}><Card><Statistic title="TillList 设备" value={stats?.till_count} prefix={<LaptopOutlined />} /></Card></Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={6}><Card><Statistic title="Var 变量" value={stats?.var_count} prefix={<CodeOutlined />} /></Card></Col>
        <Col span={6}><Card><Statistic title="Group 分组" value={stats?.group_count} prefix={<FolderOutlined />} /></Card></Col>
        <Col span={6}><Card><Statistic title="用户数" value={stats?.user_count} prefix={<UserOutlined />} /></Card></Col>
      </Row>
    </div>
  );
}
