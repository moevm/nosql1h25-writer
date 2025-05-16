import React from 'react';
import { createRoute } from '@tanstack/react-router';
import { Card, Button, Row, Col, Dropdown, Menu, Space, Typography, message } from 'antd';
import { api } from '../integrations/auth';
import type { RootRoute } from '@tanstack/react-router';
import { getUserIdFromToken } from '../integrations/auth';

const { Title, Text } = Typography;

function useUserProfile() {
  const [data, setData] = React.useState<any>(null);
  const [loading, setLoading] = React.useState(true);
  const userId = getUserIdFromToken();

  React.useEffect(() => {
    if (!userId) return;
    api.get(`/users/${userId}`)
      .then(res => setData(res.data))
      .catch(() => message.error('Ошибка загрузки профиля'))
      .finally(() => setLoading(false));
  }, [userId]);

  return { data, loading };
}

const profileMenu = (
  <Menu>
    <Menu.Item key="client">Заказчик</Menu.Item>
    <Menu.Item key="freelancer">Исполнитель</Menu.Item>
    <Menu.Divider />
    <Menu.Item key="help">Помощь</Menu.Item>
    <Menu.Item key="logout">Выход</Menu.Item>
  </Menu>
);

const balanceMenu = (
  <Menu>
    <Menu.Item key="withdraw">Вывод</Menu.Item>
    <Menu.Item key="deposit">Пополнение</Menu.Item>
  </Menu>
);

function ProfilePage() {
  const { data, loading } = useUserProfile();

  if (loading) return <div style={{ textAlign: 'center', padding: 40 }}>Загрузка...</div>;
  if (!data) return <div style={{ textAlign: 'center', padding: 40 }}>Профиль не найден</div>;

  const { displayName, email, balance, client } = data;

  return (
    <div style={{ maxWidth: 700, margin: '32px auto', background: '#f7faff', borderRadius: 16, padding: 32 }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Dropdown overlay={profileMenu} trigger={['click']}>
            <Button type="default" style={{ fontWeight: 500 }}>Настройки профиля ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Dropdown overlay={balanceMenu} trigger={['click']}>
            <Button type="default" style={{ fontWeight: 500 }}>Баланс: <b>{balance?.toLocaleString() ?? 0} руб.</b> ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Space>
            <Button type="primary">Мои заказы</Button>
            <Button>Создать заказ</Button>
          </Space>
        </Col>
      </Row>
      <Row gutter={32}>
        <Col flex="auto">
          <Card style={{ borderRadius: 12, marginBottom: 16 }}>
            <Row align="middle" gutter={16}>
              <Col>
                <div style={{
                  width: 120,
                  height: 120,
                  borderRadius: '50%',
                  background: '#fde3cf',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: 48,
                  fontWeight: 700,
                  color: '#ff7a45',
                }}>
                  {displayName?.[0] || 'И'}
                </div>
              </Col>
              <Col flex="auto">
                <Title level={4} style={{ marginBottom: 0 }}>{displayName || 'Имя'}</Title>
                <Text type="secondary">{email}</Text>
                <div style={{ margin: '8px 0' }}>
                  <Text>Город: —</Text>
                </div>
                <div>
                  <Text>Пол: —</Text>
                </div>
                <div>
                  <Text>Дата рождения: —</Text>
                </div>
                <div style={{ marginTop: 12 }}>
                  <Button>Редактировать</Button>
                </div>
              </Col>
            </Row>
          </Card>
          <Card style={{ borderRadius: 12 }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <div>
                <Text strong>Рейтинг заказчика</Text>
                <div style={{ margin: '8px 0' }}>
                  <span style={{ color: '#faad14', fontSize: 20 }}>★</span>
                  <Text style={{ fontSize: 18, marginLeft: 8 }}>{client?.rating?.toFixed(1) ?? '—'}</Text>
                </div>
                <Text type="secondary">Завершённых заказов: {client?.completedOrders ?? '—'}</Text>
              </div>
              <Button>Показать отзывы</Button>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
}

export default (parentRoute: RootRoute) =>
  createRoute({
    path: '/profile',
    component: ProfilePage,
    getParentRoute: () => parentRoute,
  }); 