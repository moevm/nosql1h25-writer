import { Button, Card, Col, Dropdown, Form, InputNumber, Modal, Row, Space, Typography, message } from 'antd';
import { Link, useNavigate } from '@tanstack/react-router';
import { useLayoutEffect, useState } from 'react';
import { useUserProfile } from '../hooks/useUserProfile';
import { roleUtils } from '../utils/role';
import { api } from '../integrations/api';
import { getUserIdFromToken } from '../integrations/auth';
import type { UserRole } from '../utils/role';
import type { MenuProps } from 'antd';

const { Title, Text } = Typography;

export default function ProfilePage() {
  const userId = getUserIdFromToken();
  const [selectedRole] = useState<UserRole>(roleUtils.getRole());
  const { data, isLoading, refetch } = useUserProfile(userId || '', selectedRole);
  const navigate = useNavigate();
  const [depositModalVisible, setDepositModalVisible] = useState(false);
  const [withdrawModalVisible, setWithdrawModalVisible] = useState(false);
  const [depositForm] = Form.useForm();
  const [withdrawForm] = Form.useForm();

  useLayoutEffect(() => {
    const originalStyle = window.getComputedStyle(document.body).overflow;
    document.body.style.overflow = depositModalVisible || withdrawModalVisible ? 'hidden' : originalStyle;
    return () => {
      document.body.style.overflow = originalStyle;
    };
  }, [depositModalVisible, withdrawModalVisible]);

  const handleDeposit = async (values: { amount: number }) => {
    try {
      await api.post('/balance/deposit', { amount: values.amount });
      message.success('Баланс успешно пополнен');
      setDepositModalVisible(false);
      depositForm.resetFields();
      await refetch();
    } catch (error) {
      message.error('Ошибка при пополнении баланса');
    }
  };

  const handleWithdraw = async (values: { amount: number }) => {
    try {
      await api.post('/balance/withdraw', { amount: values.amount });
      message.success('Средства успешно выведены');
      setWithdrawModalVisible(false);
      withdrawForm.resetFields();
      await refetch();
    } catch (error) {
      message.error('Ошибка при выводе средств');
    }
  };

  const handleRoleChange = (role: UserRole) => {
    roleUtils.setRole(role);
    window.location.reload();
  };

  const handleMyOrdersClick = () => {
    navigate({ to: '/profile/orders' });
  };

  const handleMyResponsesClick = () => {
    navigate({ to: '/profile/responses' });
  };

  const handleCreateOrderClick = () => {
    navigate({ to: '/orders/create' });
  };

  const handleOrdersClick = () => {
    navigate({ to: '/orders' });
  };

  const profileMenuItems: MenuProps['items'] = [
    {
      key: '1',
      label: 'Редактировать профиль',
      onClick: () => navigate({ to: '/profile/edit' })
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
      label: 'Пополнить баланс',
      onClick: () => setDepositModalVisible(true)
    },
    {
      key: '2',
      label: 'Вывести средства',
      onClick: () => setWithdrawModalVisible(true)
    }
  ];

  if (isLoading) return <div style={{textAlign: 'center', padding: 40}}>Загрузка...</div>;
  if (!data) return <div style={{textAlign: 'center', padding: 40}}>Профиль не найден</div>;

  const { displayName, email, balance} = data;

  return (
    <div style={{
      maxWidth: 700,
      margin: '32px auto',
      background: '#f7faff',
      borderRadius: 16,
      padding: 32,
      position: 'fixed',
      top: 32,
      left: '50%',
      transform: 'translateX(-50%)',
      width: '100%',
      maxHeight: 'calc(100vh - 64px)',
      overflow: 'auto'
    }}>
      <Row justify="space-between" align="middle" style={{marginBottom: 24}}>
        <Col>
          <Dropdown menu={{items: profileMenuItems}} trigger={['click']}>
            <Button type="default" style={{fontWeight: 500}}>Настройки профиля ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Dropdown menu={{items: balanceMenuItems}} trigger={['click']}>
            <Button type="default"
                    style={{fontWeight: 500}}>Баланс: <b>{balance.toLocaleString() || '0'} ₽</b> ▼</Button>
          </Dropdown>
        </Col>
        <Col>
          <Space>
            {selectedRole === 'client' ? (
              <>
                <Button type="primary" onClick={handleMyOrdersClick}>
                  Мои заказы
                </Button>
                <Button type="primary" onClick={handleCreateOrderClick}>
                  Создать заказ
                </Button>
              </>
            ) : (
              <>
                <Button type="primary" onClick={handleMyResponsesClick}>
                  Мои отклики
                </Button>
                <Button type="primary" onClick={handleOrdersClick}>
                  Заказы
                </Button>
              </>
            )}
          </Space>
        </Col>
      </Row>
      
      <Row>
        <Col span={24}>
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
                  <Text>
                    {selectedRole === 'client'
                    ? data.client?.description || 'Этот пользователь не установил описание'
                    : data.freelancer?.description || 'Этот пользователь не установил описание'}
                  </Text>
                </div>
                <div style={{marginTop: 12}}>
                  <Link to="/profile/edit">
                    <Button>Редактировать</Button>
                  </Link>
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
                      ? data.client?.rating.toFixed(1) ?? '—'
                      : data.freelancer?.rating.toFixed(1) ?? '—'}
                  </Text>
                </div>
              </div>
              <Button>Показать отзывы</Button>
            </div>
          </Card>
        </Col>
      </Row>

      <Modal
        title="Пополнить баланс"
        open={depositModalVisible}
        onCancel={() => {
          setDepositModalVisible(false);
          depositForm.resetFields();
        }}
        footer={null}
        maskClosable={false}
        style={{ top: 20 }}
        wrapClassName="custom-modal"
      >
        <Form form={depositForm} onFinish={handleDeposit} layout="vertical">
          <Form.Item
            name="amount"
            label="Сумма пополнения"
            rules={[
              { required: true, message: 'Введите сумму' },
              { type: 'number', min: 1, message: 'Сумма должна быть больше 0' }
            ]}
          >
            <InputNumber
              style={{ width: '100%' }}
              placeholder="Введите сумму"
              min={1}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Пополнить
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="Вывести средства"
        open={withdrawModalVisible}
        onCancel={() => {
          setWithdrawModalVisible(false);
          withdrawForm.resetFields();
        }}
        footer={null}
        maskClosable={false}
        style={{ top: 20 }}
        wrapClassName="custom-modal"
      >
        <Form form={withdrawForm} onFinish={handleWithdraw} layout="vertical">
          <Form.Item
            name="amount"
            label="Сумма вывода"
            rules={[
              { required: true, message: 'Введите сумму' },
              { type: 'number', min: 1, message: 'Сумма должна быть больше 0' },
              () => ({
                validator(_, value) {
                  if (!value || value <= data.balance) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('Недостаточно средств'));
                },
              }),
            ]}
          >
            <InputNumber
              style={{ width: '100%' }}
              placeholder="Введите сумму"
              min={1}
              max={data.balance}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Вывести
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}