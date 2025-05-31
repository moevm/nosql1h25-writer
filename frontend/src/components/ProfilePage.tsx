import { Button, Card, Col, Dropdown, Menu, Row, Space, Typography } from 'antd';
import { useUserProfile } from '../hooks/useUserProfile';
import { roleUtils, type UserRole } from '../utils/role';
import { useNavigate } from '@tanstack/react-router'
import { useState, useEffect } from 'react';

const { Title, Text } = Typography;

const balanceMenu = (
  <Menu>
    <Menu.Item key="1">Пополнить баланс</Menu.Item>
    <Menu.Item key="2">История операций</Menu.Item>
  </Menu>
);

export default function ProfilePage() {
  const { data, isLoading } = useUserProfile();
  const navigate = useNavigate();
  const [currentRole, setCurrentRole] = useState<UserRole>(roleUtils.getRole());

  const handleRoleChange = (role: UserRole) => {
    if (role === currentRole) return;
    roleUtils.setRole(role);
    navigate({ to: '/profile' });
    setCurrentRole(role);
  };

  useEffect(() => {
    const interval = setInterval(() => {
      const roleInStorage = roleUtils.getRole();
      if (roleInStorage !== currentRole) {
        setCurrentRole(roleInStorage);
      }
    }, 500);

    return () => clearInterval(interval);
  }, [currentRole]);

  const profileMenuItems = [
    {
      key: 'client',
      label: 'Заказчик',
      onClick: () => handleRoleChange('client'),
      style: currentRole === 'client' ? { color: '#1890ff', fontWeight: 'bold' } : {}
    },
    {
      key: 'freelancer',
      label: 'Исполнитель',
      onClick: () => handleRoleChange('freelancer'),
      style: currentRole === 'freelancer' ? { color: '#1890ff', fontWeight: 'bold' } : {}
    },
    { key: 'edit', label: 'Редактировать профиль' },
    { key: 'security', label: 'Настройки безопасности' },
    { key: 'notifications', label: 'Уведомления' }
  ];

  if (isLoading) return <div style={{ textAlign: 'center', padding: 40 }}>Загрузка...</div>;
  if (!data) return <div style={{ textAlign: 'center', padding: 40 }}>Профиль не найден</div>;

  const { displayName, email, balance, client } = data;

  return (
    <div style={{ maxWidth: 700, margin: '32px auto', background: '#f7faff', borderRadius: 16, padding: 32 }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Dropdown menu={{ items: profileMenuItems }} trigger={['click']}>
            <Button type="default" style={{ fontWeight: 500 }}>Настройки профиля ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Dropdown menu={{ items: balanceMenu.props.children }} trigger={['click']}>
            <Button type="default" style={{ fontWeight: 500 }}>Баланс: <b>{balance.toLocaleString() || '0'} руб.</b> ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Space>
            <Button type="primary">Мои заказы</Button>
            {currentRole === 'client' ?
              <Button>Создать заказ</Button> :
              <Button onClick={() => navigate({ to: '/orders' })}>На главную</Button>
            }
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
                  {displayName[0] || 'И'}
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
                <Text strong>
                  {currentRole === 'client' ? 'Рейтинг заказчика' : 'Рейтинг исполнителя'}
                </Text>
                <div style={{ margin: '8px 0' }}>
                  <span style={{ color: '#faad14', fontSize: 20 }}>★</span>
                  <Text style={{ fontSize: 18, marginLeft: 8 }}>
                    {currentRole === 'client'
                      ? client?.rating?.toFixed(1) ?? '—'
                      : data.freelancer?.rating?.toFixed(1) ?? '—'}
                  </Text>
                </div>
                <Text type="secondary">
                  Завершённых заказов: {currentRole === 'client'
                    ? client?.completedOrders ?? '—'
                    : data.freelancer?.completedOrders ?? '—'}
                </Text>
              </div>
              <Button>Показать отзывы</Button>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
} 