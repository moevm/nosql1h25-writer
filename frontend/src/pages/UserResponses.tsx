import { useEffect, useState } from 'react';
import { Button, Card, Col, Row, Spin, Tag, Typography } from 'antd';
import { Link, useNavigate } from '@tanstack/react-router';
import { api } from '../integrations/api';
import { formatDate } from '../utils/date';
import { getUserIdFromToken } from '../integrations/auth';

const { Title, Text } = Typography;

interface Response {
  orderId: string;
  title: string;
  coverLetter: string;
  cost: number;
  status: string;
  completionTime: number;
  createdAt: string;
}

// Порядок статусов от более раннего к более позднему
const statusOrder: Record<string, number> = {
  beginning: 0,
  negotiation: 1,
  budgeting: 2,
  work: 3,
  reviews: 4,
  finished: 5,
  dispute: 6,
};

export const UserResponses = () => {
  const userId = getUserIdFromToken();
  const navigate = useNavigate();
  const [responses, setResponses] = useState<Array<Response>>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchResponses = async () => {
      try {
        const response = await api.get(`/users/${userId}/responses`);
        // Сортируем отклики по статусу и времени создания
        const sortedResponses = response.data.responses.sort((a: Response, b: Response) => {
          // Сначала сравниваем по статусу
          const statusComparison = statusOrder[a.status] - statusOrder[b.status];
          if (statusComparison !== 0) {
            return statusComparison;
          }
          // Если статусы одинаковые, сортируем по времени создания (новые сверху)
          return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
        });
        setResponses(sortedResponses);
      } catch (error) {
        console.error('Error fetching responses:', error);
      } finally {
        setLoading(false);
      }
    };

    if (userId) {
      fetchResponses();
    }
  }, [userId]);

  const getStatusColor = (status: string) => {
    const statusColors: Record<string, string> = {
      beginning: 'blue',
      negotiation: 'orange',
      budgeting: 'purple',
      work: 'cyan',
      reviews: 'gold',
      finished: 'green',
      dispute: 'red',
    };
    return statusColors[status] || 'default';
  };

  const getStatusLabel = (status: string) => {
    const statusLabels: Record<string, string> = {
      beginning: 'Новый',
      negotiation: 'В обсуждении',
      budgeting: 'Согласование бюджета',
      work: 'В работе',
      reviews: 'На проверке',
      finished: 'Завершен',
      dispute: 'Спор',
    };
    return statusLabels[status] || status;
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div style={{ maxWidth: 1200, margin: '32px auto', padding: '0 24px' }}>
      <Title level={2} style={{ marginBottom: 24 }}>
        Мои отклики
      </Title>
      {responses.length === 0 ? (
        <Card style={{ textAlign: 'center', padding: '48px 24px' }}>
          <Title level={3} style={{ marginBottom: 16, color: '#1890ff' }}>
            У вас пока нет откликов
          </Title>
          <Text style={{ fontSize: 16, color: '#666' }}>
            Это отличное время, чтобы найти интересные проекты и предложить свою помощь! 
            Просмотрите доступные заказы и сделайте первый шаг к успешному сотрудничеству.
          </Text>
          <div style={{ marginTop: 24 }}>
            <Button type="primary" size="large" onClick={() => navigate({ to: '/orders' })}>
              Найти заказы
            </Button>
          </div>
        </Card>
      ) : (
        <Row gutter={[16, 16]}>
          {responses.map((response) => (
            <Col xs={24} key={response.orderId}>
              <Card>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 16 }}>
                  <Link 
                    to="/orders/$id" 
                    params={{ id: response.orderId }}
                    style={{ textDecoration: 'none' }}
                  >
                    <Title level={4} style={{ margin: 0, color: '#1890ff' }}>
                      {response.title}
                    </Title>
                  </Link>
                  <Tag color={getStatusColor(response.status)}>
                    {getStatusLabel(response.status)}
                  </Tag>
                </div>
                <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
                  Отклик отправлен: {formatDate(response.createdAt)}
                </Text>
                <Text style={{ display: 'block', marginBottom: 16 }}>
                  {response.coverLetter}
                </Text>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Text type="secondary">
                    Бюджет: {response.cost ? `${response.cost.toLocaleString()} ₽` : 'По договорённости'}
                  </Text>
                </div>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </div>
  );
}; 