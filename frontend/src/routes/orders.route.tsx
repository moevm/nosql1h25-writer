import React from 'react'
import { createRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Button, Card, Col, Input, Pagination, Row, Select, Spin, Slider, Space, Tag } from 'antd'
import { api } from '../integrations/auth'
import type { RootRoute } from '@tanstack/react-router'
import './orders.css'

const { Search } = Input
const { Option } = Select

// Тип заказа
interface Order {
  clientName: string
  completionTime: number
  cost: number
  description: string
  rating: number
  title: string
  status: 'new' | 'in_progress' | 'completed'
}

function OrdersList() {
  const [page, setPage] = React.useState(1)
  const [pageSize, setPageSize] = React.useState(6)
  const [search, setSearch] = React.useState('')
  const [statusFilter, setStatusFilter] = React.useState<string>('all')
  const [costRange, setCostRange] = React.useState<[number, number]>([0, 100000])
  const [sortBy, setSortBy] = React.useState<string>('newest')

  const { data, isLoading } = useQuery<{ orders: Array<Order>; total: number }>({
    queryKey: ['orders', page, pageSize, search, statusFilter, costRange, sortBy],
    queryFn: async () => {
      const params = new URLSearchParams({
        offset: String((page - 1) * pageSize),
        limit: String(pageSize),
        search,
        status: statusFilter !== 'all' ? statusFilter : '',
        minCost: String(costRange[0]),
        maxCost: String(costRange[1]),
        sortBy,
      })
      const res = await api.get(`/orders?${params}`)
      return res.data as { orders: Array<Order>; total: number }
    },
  })

  const orders: Array<Order> = data && 'orders' in data ? data.orders : []
  const total = data && 'total' in data ? data.total : 0

  const getStatusColor = (status: Order['status']) => {
    switch (status) {
      case 'new':
        return 'blue'
      case 'in_progress':
        return 'orange'
      case 'completed':
        return 'green'
      default:
        return 'default'
    }
  }

  const getStatusText = (status: Order['status']) => {
    switch (status) {
      case 'new':
        return 'Новый'
      case 'in_progress':
        return 'В работе'
      case 'completed':
        return 'Завершен'
      default:
        return status
    }
  }

  return (
    <div style={{ padding: 24 }}>
      <h1 style={{ fontSize: 24, marginBottom: 16 }}>Главная исполнителя</h1>
      <Row gutter={[16, 16]} align="middle" style={{ marginBottom: 16 }}>
        <Col>
          <Button type="primary">В профиль</Button>
        </Col>
        <Col>
          <Select 
            value={statusFilter} 
            onChange={setStatusFilter} 
            style={{ width: 150 }}
          >
            <Option value="all">Все статусы</Option>
            <Option value="new">Новые</Option>
            <Option value="in_progress">В работе</Option>
            <Option value="completed">Завершенные</Option>
          </Select>
        </Col>
        <Col>
          <Select 
            value={sortBy} 
            onChange={setSortBy} 
            style={{ width: 150 }}
          >
            <Option value="newest">Сначала новые</Option>
            <Option value="oldest">Сначала старые</Option>
            <Option value="cost_asc">По возрастанию цены</Option>
            <Option value="cost_desc">По убыванию цены</Option>
          </Select>
        </Col>
        <Col flex="auto">
          <Search
            placeholder="Найти заказ"
            onSearch={setSearch}
            allowClear
            style={{ maxWidth: 300 }}
          />
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col span={24}>
          <Card>
            <Space direction="vertical" style={{ width: '100%' }}>
              <div>Диапазон стоимости:</div>
              <Slider
                range
                value={costRange}
                onChange={(value) => setCostRange(value as [number, number])}
                min={0}
                max={100000}
                step={1000}
                marks={{
                  0: '0₽',
                  25000: '25к₽',
                  50000: '50к₽',
                  75000: '75к₽',
                  100000: '100к₽',
                }}
              />
            </Space>
          </Card>
        </Col>
      </Row>

      {isLoading ? (
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <Spin size="large" />
        </div>
      ) : (
        <Row gutter={[16, 16]}>
          {orders.map((order, idx) => (
            <Col xs={24} sm={12} md={8} key={idx}>
              <Card
                title={
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <div style={{
                      width: 48,
                      height: 48,
                      borderRadius: '50%',
                      background: '#fde3cf',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontWeight: 700,
                      fontSize: 24,
                    }}>
                      {order.clientName[0]}
                    </div>
                    <div>
                      <div style={{ fontWeight: 600 }}>{order.clientName}</div>
                      <div style={{ color: '#faad14', fontSize: 14 }}>
                        {'★'.repeat(Math.round(order.rating))}
                        <span style={{ color: '#888', marginLeft: 4 }}>{order.rating.toFixed(1)}</span>
                      </div>
                    </div>
                  </div>
                }
                extra={<Tag color={getStatusColor(order.status)}>{getStatusText(order.status)}</Tag>}
                className="order-card"
                style={{ 
                  height: '100%',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.09)',
                }}
                hoverable
              >
                <div style={{ fontWeight: 600, marginBottom: 4, fontSize: 16 }}>{order.title}</div>
                <div style={{ 
                  marginBottom: 8, 
                  minHeight: 48,
                  color: '#666',
                  lineHeight: 1.5
                }}>
                  {order.description.length > 80
                    ? order.description.slice(0, 80) + '...'
                    : order.description}
                </div>
                <div style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center',
                  borderTop: '1px solid #f0f0f0',
                  paddingTop: 12,
                  marginTop: 12
                }}>
                  <span style={{ color: '#666' }}>
                    ⏰ {Math.floor(order.completionTime / (24 * 60 * 60 * 1000))} дня {Math.floor((order.completionTime % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000))} часов
                  </span>
                  <span style={{ 
                    fontWeight: 700, 
                    fontSize: 18,
                    color: '#1890ff'
                  }}>
                    💰 {order.cost.toLocaleString()} ₽
                  </span>
                </div>
              </Card>
            </Col>
          ))}
        </Row>
      )}
      <div style={{ marginTop: 24, textAlign: 'center' }}>
        <Pagination
          current={page}
          pageSize={pageSize}
          total={total}
          onChange={setPage}
          onShowSizeChange={(_, size) => setPageSize(size)}
          showSizeChanger
          pageSizeOptions={[6, 12, 18, 24]}
          showTotal={(total) => `Всего ${total} заказов`}
        />
      </div>
    </div>
  )
}

export default (parentRoute: RootRoute) =>
  createRoute({
    path: '/orders',
    component: OrdersList,
    getParentRoute: () => parentRoute,
  }) 