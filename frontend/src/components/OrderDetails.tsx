import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { useParams } from '@tanstack/react-router'
import { Avatar, Button, Card, Input, List, Spin, Tag } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { api } from '../integrations/api'
import { roleUtils } from '../utils/role'
import { formatCompletionTime } from '../utils/time'

interface OrderDetailsType {
  order: {
    id: string
    clientId: string
    clientName: string
    clientRating: number
    completionTime: number
    cost: number
    description: string
    title: string
    status: string
    createdAt: string
    updatedAt: string
    responses?: Array<{
      freelancerName: string
      freelancerId: string
      coverLetter: string
      createdAt: string
    }>
  }
  isClient: boolean
  isFreelancer: boolean
  hasActiveResponse: boolean
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'beginning': return 'blue'
    case 'negotiation': return 'orange'
    case 'budgeting': return 'purple'
    case 'work': return 'cyan'
    case 'reviews': return 'gold'
    case 'finished': return 'green'
    case 'dispute': return 'red'
    default: return 'default'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'beginning': return 'Новый'
    case 'negotiation': return 'В обсуждении'
    case 'budgeting': return 'Согласование бюджета'
    case 'work': return 'В работе'
    case 'reviews': return 'На проверке'
    case 'finished': return 'Завершен'
    case 'dispute': return 'Спор'
    default: return status
  }
}

const formatRating = (rating: number) => {
  if (rating === 0) {
    return 'Нет отзывов'
  }
  return (
    <>
      {'★'.repeat(Math.round(rating))} 
      <span style={{ marginLeft: 6, color: '#888' }}>{rating.toFixed(1)}</span>
    </>
  )
}

const OrderDetails: React.FC = () => {
  const { id } = useParams({ from: '/orders/$id' })
  const { data, isLoading } = useQuery<OrderDetailsType>({
    queryKey: ['order', id],
    queryFn: () => api.get(`/orders/${id}`).then(res => res.data)
  })

  if (isLoading) {
    return (
      <div style={{ 
        textAlign: 'center', 
        padding: '80px 0',
        minHeight: 'calc(100vh - 64px)'
      }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!data) {
    return (
      <div style={{ 
        textAlign: 'center', 
        padding: '80px 0',
        minHeight: 'calc(100vh - 64px)'
      }}>
        Заказ не найден
      </div>
    )
  }

  const { order } = data
  const isCurrentRoleClient = roleUtils.getRole() === 'client'

  return (
    <div style={{ 
      maxWidth: 800, 
      margin: '0 auto',
      padding: '80px 16px',
      minHeight: 'calc(100vh - 64px)'
    }}>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
            <div style={{
              width: 64,
              height: 64,
              borderRadius: '50%',
              backgroundColor: '#fde3cf',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 700,
              fontSize: 28,
            }}>
              {order.clientName[0]}
            </div>
            <div>
              <div style={{ fontWeight: 600 }}>{order.clientName}</div>
              <div style={{ color: '#faad14' }}>
                {formatRating(order.clientRating)}
              </div>
            </div>
          </div>
          <div style={{ textAlign: 'right' }}>
            <div style={{ marginBottom: 8 }}>
              ⏰ {formatCompletionTime(order.completionTime)}
            </div>
            <div style={{ fontSize: 18, fontWeight: 700, color: '#1890ff' }}>
              💰 {order.cost ? `${order.cost.toLocaleString()} ₽` : 'По договорённости'}
            </div>
          </div>
        </div>

        <Tag color={getStatusColor(order.status)}>{getStatusText(order.status)}</Tag>

        <h2 style={{ marginTop: 24 }}>{order.title}</h2>

        <div style={{ 
          marginTop: 24,
          whiteSpace: 'pre-wrap',
          lineHeight: 1.6,
          color: '#333'
        }}>
          {order.description}
        </div>

        {data.isClient && isCurrentRoleClient && (
          <div style={{ 
            marginTop: 32,
            backgroundColor: '#fafafa',
            borderRadius: 8,
            padding: 16
          }}>
            <h3 style={{ 
              fontSize: 18,
              marginBottom: 16,
              color: '#1890ff'
            }}>
              Отклики {order.responses && order.responses.length > 0 ? `(${order.responses.length})` : ''}
            </h3>
            
            {order.responses && order.responses.length > 0 ? (
              <List
                dataSource={order.responses}
                renderItem={(response) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={<Avatar icon={<UserOutlined />} />}
                      title={response.freelancerName}
                      description={
                        <div>
                          <div style={{ marginBottom: 8 }}>{response.coverLetter}</div>
                          <div style={{ color: '#888', fontSize: 12 }}>
                            {new Date(response.createdAt).toLocaleString('ru-RU')}
                          </div>
                        </div>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <div style={{ 
                textAlign: 'center',
                color: '#666',
                padding: '24px 0'
              }}>
                Откликов пока нет. Возможно, стоит немного подождать или уточнить требования к заказу.
              </div>
            )}
          </div>
        )}

        {!data.isClient && !data.hasActiveResponse && (
          <div style={{ marginTop: 32 }}>
            <Input.TextArea placeholder="Написать заказчику..." rows={4} style={{ marginBottom: 12 }} />
            <Button type="primary">Готов взяться</Button>
          </div>
        )}
      </Card>
    </div>
  )
}

export default OrderDetails;