import React from 'react'
import { createRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Card, Row, Col, Pagination, Button, Input, Select, Spin } from 'antd'
import type { RootRoute } from '@tanstack/react-router'
import { api } from '../integrations/auth'

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
}

function OrdersList() {
  const [page, setPage] = React.useState(1)
  const [pageSize, setPageSize] = React.useState(6)
  const [search, setSearch] = React.useState('')

  const { data, isLoading } = useQuery<{ orders: Order[]; total: number }>({
    queryKey: ['orders', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams({
        offset: String((page - 1) * pageSize),
        limit: String(pageSize),
      })
      const res = await api.get(`/orders?${params}`)
      return res.data as { orders: Order[]; total: number }
    },
  })

  const orders: Order[] = data && 'orders' in data ? data.orders : []
  const total = data && 'total' in data ? data.total : 0

  return (
    <div style={{ padding: 24 }}>
      <h1 style={{ fontSize: 24, marginBottom: 16 }}>Главная исполнителя</h1>
      <Row gutter={[16, 16]} align="middle" style={{ marginBottom: 16 }}>
        <Col>
          <Button type="primary">В профиль</Button>
        </Col>
        <Col>
          <Select defaultValue="all" style={{ width: 120 }}>
            <Option value="all">Фильтры</Option>
            {/* Добавить опции фильтрации */}
          </Select>
        </Col>
        <Col>
          <Select defaultValue="default" style={{ width: 150 }}>
            <Option value="default">Сортировать</Option>
            {/* Добавить опции сортировки */}
          </Select>
        </Col>
        <Col flex="auto">
          <Search
            placeholder="Найти вакансию"
            onSearch={setSearch}
            allowClear
            style={{ maxWidth: 300 }}
          />
        </Col>
      </Row>
      {isLoading ? (
        <Spin size="large" />
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
                bordered={false}
                style={{ background: '#e6f4ff' }}
              >
                <div style={{ fontWeight: 600, marginBottom: 4 }}>{order.title}</div>
                <div style={{ marginBottom: 8, minHeight: 48 }}>
                  {order.description.length > 80
                    ? order.description.slice(0, 80) + '...'
                    : order.description}
                </div>
                <div style={{ display: 'flex', gap: 16, alignItems: 'center', marginBottom: 4 }}>
                  <span>⏰ {Math.floor(order.completionTime / (24 * 60 * 60 * 1000))} дня {Math.floor((order.completionTime % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000))} часов</span>
                  <span style={{ fontWeight: 700, fontSize: 18 }}>💰 {order.cost} руб</span>
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