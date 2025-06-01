import { useEffect, useState } from 'react';
import { Card, Col, Row, Spin, Tag, Typography } from 'antd';
import { Link } from '@tanstack/react-router';
import { api } from '../integrations/api';
import { formatDate } from '../utils/date';
import { getUserIdFromToken } from '../integrations/auth';

const { Title, Text } = Typography;

interface Order {
  id: string;
  title: string;
  description: string;
  cost: number;
  status: string;
  completionTime: number;
  createdAt: string;
  updatedAt: string;
  clientId: string;
  freelancerId: string | null;
  totalResponses: number;
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

export const UserOrders = () => {
  const userId = getUserIdFromToken();
  const [orders, setOrders] = useState<Array<Order>>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const response = await api.get(`/users/${userId}/orders`);
        // Сортируем заказы по статусу и времени создания
        const sortedOrders = response.data.orders.sort((a: Order, b: Order) => {
          // Сначала сравниваем по статусу
          const statusComparison = statusOrder[a.status] - statusOrder[b.status];
          if (statusComparison !== 0) {
            return statusComparison;
          }
          // Если статусы одинаковые, сортируем по времени создания (новые сверху)
          return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
        });
        setOrders(sortedOrders);
      } catch (error) {
        console.error('Error fetching orders:', error);
      } finally {
        setLoading(false);
      }
    };

    if (userId) {
      fetchOrders();
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
        Мои заказы
      </Title>
      <Row gutter={[16, 16]}>
        {orders.map((order) => (
          <Col xs={24} key={order.id}>
            <Card>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 16 }}>
                <Link 
                  to="/orders/$id" 
                  params={{ id: order.id }}
                  style={{ textDecoration: 'none' }}
                >
                  <Title level={4} style={{ margin: 0, color: '#1890ff' }}>
                    {order.title}
                  </Title>
                </Link>
                <Tag color={getStatusColor(order.status)}>
                  {getStatusLabel(order.status)}
                </Tag>
              </div>
              <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
                Создан: {formatDate(order.createdAt)}
              </Text>
              <Text style={{ display: 'block', marginBottom: 16 }}>
                {order.description}
              </Text>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Text type="secondary">
                  Бюджет: {order.cost ? `${order.cost.toLocaleString()} ₽` : 'По договорённости'}
                </Text>
                <Text type="secondary">
                  Откликов: {order.totalResponses || 0}
                </Text>
              </div>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}; 