import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { useParams } from '@tanstack/react-router'
import { Card, Input, Button, Spin, Tag } from 'antd'
import { api } from '../integrations/auth'

interface OrderDetailsType {
  id: string
  clientName: string
  completionTime: number
  cost: number
  description: string
  rating: number
  title: string
  status: 'new' | 'in_progress' | 'completed'
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'new': return 'blue'
    case 'in_progress': return 'orange'
    case 'completed': return 'green'
    default: return 'default'
  }
}

const OrderDetails: React.FC = () => {
  const { id } = useParams({ strict: false }) as { id: string }

  const { data, isLoading } = useQuery<OrderDetailsType>({
    queryKey: ['order', id],
    queryFn: async () => {
      const res = await api.get(`/orders/${id}`)
      return res.data
    },
  })

  if (isLoading || !data) {
    return <div style={{ textAlign: 'center', padding: 50 }}><Spin size="large" /></div>
  }

  return (
    <div style={{ padding: 24 }}>
      <h1 style={{ fontSize: 24, marginBottom: 16 }}>Подробности заказа</h1>
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
              {data.clientName[0]}
            </div>
            <div>
              <div style={{ fontWeight: 600 }}>{data.clientName}</div>
              <div style={{ color: '#faad14' }}>
                {'★'.repeat(Math.round(data.rating))} 
                <span style={{ marginLeft: 6, color: '#888' }}>{data.rating.toFixed(1)}</span>
              </div>
            </div>
          </div>
          <div style={{ textAlign: 'right' }}>
            <div style={{ marginBottom: 8 }}>
              ⏰ {Math.floor(data.completionTime / (24 * 60 * 60 * 1000))} д.{' '}
              {Math.floor((data.completionTime % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000))} ч.
            </div>
            <div style={{ fontSize: 18, fontWeight: 700, color: '#1890ff' }}>
              💰 {data.cost.toLocaleString()} ₽
            </div>
          </div>
        </div>

        <Tag color={getStatusColor(data.status)}>{data.status}</Tag>

        <h2 style={{ marginTop: 24 }}>{data.title}</h2>
        <p style={{ lineHeight: 1.6 }}>{data.description}</p>

        <div style={{ marginTop: 32 }}>
          <Input.TextArea placeholder="Написать заказчику..." rows={4} style={{ marginBottom: 12 }} />
          <Button type="primary">Готов взяться</Button>
        </div>
      </Card>
    </div>
  )
}

export default OrderDetails
