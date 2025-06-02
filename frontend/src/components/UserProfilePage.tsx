import { Button, Card, Col, Radio, Row, Typography } from 'antd';
import { useState } from 'react';
import { useUserProfile } from '../hooks/useUserProfile';
import type { UserRole } from '../utils/role';

const { Title, Text } = Typography;

interface UserProfilePageProps {
  userId: string;
}

export default function UserProfilePage({ userId }: UserProfilePageProps) {
  const [selectedRole, setSelectedRole] = useState<UserRole>('client');
  const { data, isLoading } = useUserProfile(userId, selectedRole);

  if (isLoading) return <div style={{textAlign: 'center', padding: 40}}>Загрузка...</div>;
  if (!data) return <div style={{textAlign: 'center', padding: 40}}>Профиль не найден</div>;

  const { displayName } = data;

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
          <Radio.Group value={selectedRole} onChange={(e) => setSelectedRole(e.target.value)}>
            <Radio.Button value="client">Заказчик</Radio.Button>
            <Radio.Button value="freelancer">Исполнитель</Radio.Button>
          </Radio.Group>
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
                <div style={{margin: '8px 0'}}>
                  <Text>
                    {selectedRole === 'client'
                      ? data.client?.description || 'Этот пользователь не установил описание'
                      : data.freelancer?.description || 'Этот пользователь не установил описание'}
                  </Text>
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
    </div>
  );
} 