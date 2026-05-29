import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card, Descriptions, Tag, Button, Space, Spin, message, Modal, Timeline, Tabs, Empty, Typography } from 'antd';
import { ArrowLeftOutlined, HistoryOutlined, CheckCircleOutlined, CloseCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { getTillList, listSyncReports } from '../api/till_list';
import type { TillList, SyncReport, ModuleExec, ModuleStep } from '../api/till_list';
import { listAllGroups } from '../api/group';

const { Text } = Typography;

const envColor: Record<string, string> = { pr: 'red', pp: 'orange', lab: 'blue' };
const reportStatusColor: Record<number, string> = { 0: 'green', 1: 'red' };
const reportStatusIcon: Record<number, any> = { 0: <CheckCircleOutlined />, 1: <CloseCircleOutlined /> };

function parseModules(raw: string): ModuleExec[] {
  try { return JSON.parse(raw); } catch { return []; }
}

function parseFact(body: string): Array<{ Key: string; Value: string; Type: number }> {
  try {
    const parsed = JSON.parse(body);
    return parsed.Fact || [];
  } catch { return []; }
}

function formatTime(t: string): string {
  if (!t) return '-';
  return t.replace('T', ' ').substring(0, 19);
}

function StepsTable({ steps }: { steps: ModuleStep[] }) {
  if (!steps || steps.length === 0) return <Text type="secondary">无步骤</Text>;
  return (
    <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 12 }}>
      <thead>
        <tr style={{ background: '#fafafa' }}>
          <th style={{ padding: '4px 8px', textAlign: 'left', border: '1px solid #e8e8e8', width: 110 }}>步骤</th>
          <th style={{ padding: '4px 8px', textAlign: 'center', border: '1px solid #e8e8e8', width: 60 }}>结果</th>
          <th style={{ padding: '4px 8px', textAlign: 'center', border: '1px solid #e8e8e8', width: 50 }}>耗时</th>
          <th style={{ padding: '4px 8px', textAlign: 'left', border: '1px solid #e8e8e8' }}>输出</th>
        </tr>
      </thead>
      <tbody>
        {steps.map((step, si) => (
          <tr key={si}>
            <td style={{ padding: '4px 8px', border: '1px solid #e8e8e8', fontFamily: 'monospace', fontWeight: 600 }}>
              {step.Name}
            </td>
            <td style={{ padding: '4px 8px', border: '1px solid #e8e8e8', textAlign: 'center' }}>
              <Tag color={step.Status ? 'green' : 'red'} style={{ fontSize: 11 }}>
                {step.Status ? '成功' : '失败'}
              </Tag>
            </td>
            <td style={{ padding: '4px 8px', border: '1px solid #e8e8e8', textAlign: 'center', color: '#666' }}>
              {step.Duration}s
            </td>
            <td style={{ padding: '4px 8px', border: '1px solid #e8e8e8', fontFamily: 'monospace', fontSize: 11, whiteSpace: 'pre-wrap', wordBreak: 'break-all', maxWidth: 500 }}>
              {step.Output || '-'}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function ModulePanel({ mod }: { mod: ModuleExec }) {
  return (
    <div style={{ marginBottom: 8, border: '1px solid #e8e8e8', borderRadius: 6, overflow: 'hidden' }}>
      <div style={{
        display: 'flex', alignItems: 'center', gap: 8, padding: '8px 12px',
        background: mod.Status ? '#f6ffed' : '#fff2f0', borderBottom: '1px solid #e8e8e8',
      }}>
        <Tag color={mod.Status ? 'green' : 'red'} style={{ marginRight: 0 }}>
          {mod.Status ? '成功' : '失败'}
        </Tag>
        <span style={{ fontWeight: 600, fontSize: 13 }}>{mod.Name}</span>
        <span style={{ color: '#999', fontSize: 12, marginLeft: 'auto' }}>
          <ClockCircleOutlined style={{ marginRight: 4 }} />{mod.Duration}s
        </span>
      </div>
      <div style={{ padding: 8 }}>
        <StepsTable steps={mod.Steps} />
      </div>
    </div>
  );
}

function ReportDetail({ report }: { report: SyncReport }) {
  const modules = parseModules(report.module_execution);
  const totalSteps = modules.reduce((acc, m) => acc + (m.Steps?.length || 0), 0);
  const passedSteps = modules.reduce((acc, m) => acc + (m.Steps?.filter(s => s.Status)?.length || 0), 0);

  return (
    <div>
      <Descriptions column={4} size="small" style={{ marginBottom: 12 }}>
        <Descriptions.Item label="同步状态">
          <Tag icon={reportStatusIcon[report.status]} color={reportStatusColor[report.status]}>
            {report.status === 0 ? '成功' : '失败'}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="总耗时">{report.duration}s</Descriptions.Item>
        <Descriptions.Item label="模块数">{modules.length}</Descriptions.Item>
        <Descriptions.Item label="步骤进度">
          <Text style={{ fontSize: 13 }}>
            <Text type="success">{passedSteps}</Text>
            <Text type="secondary"> / {totalSteps}</Text>
          </Text>
        </Descriptions.Item>
      </Descriptions>
      {modules.length === 0 ? (
        <Empty description="无模块执行数据" image={Empty.PRESENTED_IMAGE_SIMPLE} />
      ) : (
        modules.map((mod, mi) => <ModulePanel key={mi} mod={mod} />)
      )}
    </div>
  );
}

export default function TillListDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [till, setTill] = useState<TillList | null>(null);
  const [reports, setReports] = useState<SyncReport[]>([]);
  const [selectedReport, setSelectedReport] = useState<SyncReport | null>(null);
  const [showBody, setShowBody] = useState(false);
  const [groupsMap, setGroupsMap] = useState<Record<number, string>>({});

  useEffect(() => {
    listAllGroups().then(list => {
      const map: Record<number, string> = {};
      list.forEach(g => { map[g.id] = g.name; });
      setGroupsMap(map);
    }).catch(() => {});
  }, []);

  useEffect(() => {
    if (!id || id === 'new') return;
    setLoading(true);
    Promise.all([
      getTillList(Number(id)),
      listSyncReports(Number(id)),
    ])
      .then(([t, r]) => {
        setTill(t);
        setReports(r);
        if (r.length > 0) setSelectedReport(r[0]);
      })
      .catch(() => message.error('加载失败'))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <Spin style={{ display: 'block', marginTop: 120 }} />;
  if (!till) return null;

  const facts = parseFact(till.request_body);
  const groupName = till.group_id ? (groupsMap[till.group_id] || `${till.group_id}`) : '-';

  return (
    <div style={{ maxWidth: 1100, margin: '0 auto' }}>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/till-lists')}>返回列表</Button>
        </Space>
        <Button icon={<HistoryOutlined />} onClick={() => setShowBody(true)}>查看原始上报数据</Button>
      </div>

      <div style={{
        marginBottom: 16, padding: '12px 16px', background: '#f5f7fa',
        borderRadius: 8, border: '1px solid #e8ecf0',
        display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '8px 24px', fontSize: 13,
      }}>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>设备名称</span><Text strong>{till.host_name}</Text></span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>MAC</span><Text code style={{ fontSize: 12 }}>{till.mac_address || '-'}</Text></span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>IP</span><Text code style={{ fontSize: 12 }}>{till.ip || '-'}</Text></span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>环境</span>{till.env ? <Tag color={envColor[till.env]} style={{ marginRight: 0 }}>{till.env}</Tag> : '-'}</span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>地区</span>{till.location || '-'}</span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>门店号</span>{till.store_number || '-'}</span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>硬件型号</span>{till.hardware_model || '-'}</span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>分组</span>{groupName}</span>
        <span><span style={{ color: '#8c8c8c', marginRight: 6 }}>最近活跃</span>{reports.length > 0 ? formatTime(reports[0].created_at) : (till.last_seen || '-')}</span>
      </div>

      <Tabs defaultActiveKey="info" items={[
        {
          key: 'info',
          label: <span><CheckCircleOutlined /> 设备信息</span>,
          children: (
            <>

              {facts.length > 0 && (
                <Card size="small" title="设备属性 (Fact)" style={{ marginBottom: 16 }}>
                  <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 13 }}>
                    <thead>
                      <tr style={{ background: '#f5f5f5' }}>
                        <th style={{ padding: '6px 12px', textAlign: 'left', border: '1px solid #e8e8e8' }}>Key</th>
                        <th style={{ padding: '6px 12px', textAlign: 'left', border: '1px solid #e8e8e8' }}>Value</th>
                        <th style={{ padding: '6px 12px', textAlign: 'center', border: '1px solid #e8e8e8', width: 80 }}>类型</th>
                      </tr>
                    </thead>
                    <tbody>
                      {facts.map((f, i) => (
                        <tr key={i}>
                          <td style={{ padding: '6px 12px', border: '1px solid #e8e8e8', fontFamily: 'monospace', fontSize: 12 }}>{f.Key}</td>
                          <td style={{ padding: '6px 12px', border: '1px solid #e8e8e8', fontFamily: 'monospace', fontSize: 12, maxWidth: 400, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }} title={f.Value}>{f.Value}</td>
                          <td style={{ padding: '6px 12px', border: '1px solid #e8e8e8', textAlign: 'center' }}>
                            <Tag color={f.Type === 1 ? 'blue' : 'default'}>{f.Type === 1 ? '版本' : '静态'}</Tag>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </Card>
              )}
            </>
          ),
        },
        {
          key: 'reports',
          label: <span><HistoryOutlined /> 同步报告 <Tag style={{ marginLeft: 4 }}>{reports.length}</Tag></span>,
          children: (
            <div style={{ display: 'flex', gap: 16 }}>
              <div style={{ width: 220, flexShrink: 0 }}>
                <Card size="small" title={`同步历史 (${reports.length})`} style={{ position: 'sticky', top: 16 }}>
                  {reports.length === 0 ? (
                    <Empty description="暂无同步记录" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                  ) : (
                    <div style={{ maxHeight: 560, overflowY: 'auto', paddingRight: 4 }}>
                    <Timeline
                      items={reports.map((r) => ({
                        color: r.status === 0 ? 'green' : 'red',
                        dot: r.status === 0 ? <CheckCircleOutlined style={{ color: '#52c41a' }} /> : <CloseCircleOutlined style={{ color: '#ff4d4f' }} />,
                        children: (
                          <div
                            onClick={() => setSelectedReport(r)}
                            style={{
                              cursor: 'pointer',
                              padding: '4px 8px',
                              borderRadius: 4,
                              background: selectedReport?.id === r.id ? '#e6f4ff' : 'transparent',
                              transition: 'background 0.2s',
                            }}
                          >
                            <div style={{ fontWeight: selectedReport?.id === r.id ? 600 : 400, fontSize: 13 }}>
                              {formatTime(r.created_at)}
                            </div>
                            <Tag color={r.status === 0 ? 'green' : 'red'} style={{ fontSize: 10, marginTop: 2 }}>
                              {r.status === 0 ? '成功' : '失败'} · {r.duration}s
                            </Tag>
                          </div>
                        ),
                      }))}
                    />
                    </div>
                  )}
                </Card>
              </div>
              <div style={{ flex: 1, minWidth: 0 }}>
                {selectedReport ? (
                  <Card
                    size="small"
                    title={
                      <Space>
                        <HistoryOutlined />
                        <span>报告详情 — {formatTime(selectedReport.created_at)}</span>
                      </Space>
                    }
                    extra={
                      <Button size="small" icon={<HistoryOutlined />} onClick={() => setShowBody(true)}>
                        查看原始JSON
                      </Button>
                    }
                  >
                    <ReportDetail report={selectedReport} />
                  </Card>
                ) : (
                  <Card><Empty description="请从左侧选择一条同步记录" image={Empty.PRESENTED_IMAGE_SIMPLE} /></Card>
                )}
              </div>
            </div>
          ),
        },
      ]} />

      <Modal
        title="原始上报数据 (Request Body)"
        open={showBody}
        onCancel={() => setShowBody(false)}
        footer={null}
        width={900}
      >
        <pre style={{ fontSize: 11, maxHeight: 600, overflow: 'auto', background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
          {selectedReport?.request_body || till.request_body
            ? JSON.stringify(JSON.parse(selectedReport?.request_body || till.request_body || '{}'), null, 2)
            : '无上报数据'}
        </pre>
      </Modal>
    </div>
  );
}
