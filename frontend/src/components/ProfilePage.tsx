import { Button, Card, Col, Dropdown, Row, Space, Typography } from 'antd';
import { Link, useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { useUserProfile } from '../hooks/useUserProfile';
import { roleUtils } from '../utils/role';
import type { UserRole } from '../utils/role';
import type { MenuProps } from 'antd';

const { Title, Text } = Typography;

export default function ProfilePage() {
  const { data, isLoading } = useUserProfile();
  const [selectedRole] = useState<UserRole>(roleUtils.getRole());
  const navigate = useNavigate();

  const handleRoleChange = (role: UserRole) => {
    roleUtils.setRole(role);
    window.location.reload();
  };

  const handleMyOrdersClick = () => {
    navigate({ to: '/profile/orders' });
  };

  const profileMenuItems: MenuProps['items'] = [
    {
      key: '1',
      label: 'Редактировать профиль'
    },
    {
      type: 'divider' as const
    },
    {
      key: 'role',
      label: 'Роль на сайте',
      children: [
        {
          key: 'client',
          label: 'Заказчик',
          onClick: () => handleRoleChange('client'),
          style: {
            backgroundColor: selectedRole === 'client' ? '#e6f7ff' : undefined,
            color: selectedRole === 'client' ? '#1890ff' : undefined
          }
        },
        {
          key: 'freelancer',
          label: 'Исполнитель',
          onClick: () => handleRoleChange('freelancer'),
          style: {
            backgroundColor: selectedRole === 'freelancer' ? '#e6f7ff' : undefined,
            color: selectedRole === 'freelancer' ? '#1890ff' : undefined
          }
        }
      ]
    }
  ];

  const balanceMenuItems: MenuProps['items'] = [
    {
      key: '1',
      label: 'Пополнить баланс'
    },
    {
      key: '2',
      label: 'История операций'
    }
  ];

  if (isLoading) return <div style={{textAlign: 'center', padding: 40}}>Загрузка...</div>;
  if (!data) return <div style={{textAlign: 'center', padding: 40}}>Профиль не найден</div>;

  const { displayName, email, balance, client } = data;

  return (
    <div style={{maxWidth: 700, margin: '32px auto', background: '#f7faff', borderRadius: 16, padding: 32}}>
      <Row justify="space-between" align="middle" style={{marginBottom: 24}}>
        <Col>
          <Dropdown menu={{items: profileMenuItems}} trigger={['click']}>
            <Button type="default" style={{fontWeight: 500}}>Настройки профиля ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Dropdown menu={{items: balanceMenuItems}} trigger={['click']}>
            <Button type="default"
                    style={{fontWeight: 500}}>Баланс: <b>{balance.toLocaleString() || '0'} руб.</b> ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Space>
            <Button type="primary" onClick={handleMyOrdersClick}>Мои заказы</Button>
            {selectedRole === 'client' ?
              <Link to="/orders/create">
                <Button>Создать заказ</Button>
              </Link> :
              <Link to="/orders">
                <Button>Заказы</Button>
              </Link>
            }
          </Space>
        </Col>
      </Row>
      
      <Col flex="auto">
        <Card style={{borderRadius: 12, marginBottom: 16}}>
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
              <Title level={4} style={{marginBottom: 0}}>{displayName || 'Имя'}</Title>
              <Text type="secondary">{email}</Text>
              <div style={{margin: '8px 0'}}>
                <Text>Город: —</Text>
              </div>
              <div>
                <Text>Пол: —</Text>
              </div>
              <div>
                <Text>Дата рождения: —</Text>
              </div>
              <div style={{marginTop: 12}}>
                <Button>Редактировать</Button>
              </div>
            </Col>
          </Row>
        </Card>
        
        <Card style={{borderRadius: 12}}>
          <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
            <div>
              <Text strong>
                {selectedRole === 'client' ? 'Рейтинг заказчика' : 'Рейтинг исполнителя'}
              </Text>
              <div style={{margin: '8px 0'}}>
                <span style={{color: '#faad14', fontSize: 20}}>★</span>
                <Text style={{fontSize: 18, marginLeft: 8}}>
                  {selectedRole === 'client'
                    ? client?.rating.toFixed(1) ?? '—'
                    : data.freelancer?.rating.toFixed(1) ?? '—'}
                </Text>
              </div>
              <Text type="secondary">
                Завершённых заказов: {selectedRole === 'client'
                  ? client?.completedOrders ?? '—'
                  : data.freelancer?.completedOrders ?? '—'}
              </Text>
            </div>
            <Button>Показать отзывы</Button>
          </div>
        </Card>
      </Col>
    </div>
  );
}