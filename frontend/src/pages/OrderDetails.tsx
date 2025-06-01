import { Card, Col, Row, Tag, Typography } from 'antd';
import { getUserIdFromToken } from '../integrations/auth';
import { formatDate } from '../utils/date';

const { Title, Text } = Typography;

interface OrderDetailsProps {
  order: {
    title: string;
    clientName: string;
    clientEmail: string;
    freelancerId: string | null;
    status: string;
    description: string;
    cost: number;
    completionTime: number;
    createdAt: string;
    updatedAt: string;
  };
  isFreelancer: boolean;
}

const getStatusColor = (status: string): string => {
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

const getStatusLabel = (status: string): string => {
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

export const OrderDetails = ({ order, isFreelancer }: OrderDetailsProps) => {
  const userId = getUserIdFromToken();

  return (
    <div style={{ maxWidth: 1200, margin: '32px auto', padding: '0 24px' }}>
      <Title level={2} style={{ marginBottom: 24 }}>
        {order.title}
      </Title>
      <Row gutter={[16, 16]}>
        <Col xs={24}>
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 16 }}>
              <div>
                <Title level={4} style={{ margin: 0 }}>
                  Заказчик: {order.clientName}
                </Title>
                <Text type="secondary" style={{ display: 'block', marginTop: 4 }}>
                  {order.clientEmail}
                </Text>
                {isFreelancer && order.freelancerId === userId && (
                  <Tag color="success" style={{ marginTop: 8 }}>
                    Вас выбрали исполнителем
                  </Tag>
                )}
              </div>
              <Tag color={getStatusColor(order.status)}>
                {getStatusLabel(order.status)}
              </Tag>
            </div>
            <Text style={{ display: 'block', marginBottom: 16 }}>
              {order.description}
            </Text>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text type="secondary" style={{ display: 'block' }}>
                  Создан: {formatDate(order.createdAt)}
                </Text>
                <Text type="secondary" style={{ display: 'block' }}>
                  Обновлен: {formatDate(order.updatedAt)}
                </Text>
              </div>
              <div>
                <Text type="secondary" style={{ display: 'block' }}>
                  Бюджет: {order.cost ? `${order.cost.toLocaleString()} ₽` : 'По договорённости'}
                </Text>
                <Text type="secondary" style={{ display: 'block' }}>
                  Срок выполнения: {formatDate(new Date(Date.now() + order.completionTime).toISOString())}
                </Text>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default OrderDetails; 