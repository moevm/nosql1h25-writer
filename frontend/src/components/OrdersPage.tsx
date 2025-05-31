import React from 'react'
import { Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import {
  Button,
  Card,
  Col,
  Input,
  InputNumber,
  Pagination,
  Row,
  Select,
  Space,
  Spin
} from 'antd'
import { CloseOutlined } from '@ant-design/icons'
import { api } from '../integrations/api'
import { formatCompletionTime } from '../utils/time'
import './orders.css'

const { Search } = Input
const { Option } = Select

// Тип заказа
interface Order {
  id: string
  clientName: string
  completionTime: number
  cost: number
  description: string
  rating: number
  title: string
  status?: 'beginning' | 'negotiation' | 'budgeting' | 'work' | 'reviews' | 'finished' | 'dispute'
  responses?: Array<{
    freelancerName: string
    freelancerId: string
    coverLetter: string
    createdAt: string
  }>
  isClient?: boolean
  isFreelancer?: boolean
  hasActiveResponse?: boolean
}

export function OrdersPage() {
  const [page, setPage] = React.useState(1)
  const [pageSize, setPageSize] = React.useState(6)
  const [search, setSearch] = React.useState('')
  const [minCost, setMinCost] = React.useState<number | null>(null)
  const [maxCost, setMaxCost] = React.useState<number | null>(null)
  const [minTime, setMinTime] = React.useState<number | null>(null)
  const [maxTime, setMaxTime] = React.useState<number | null>(null)
  const [sortBy, setSortBy] = React.useState<string>('newest')

  const { data, isLoading } = useQuery<{ orders: Array<Order>; total: number }>({
    queryKey: ['orders', page, pageSize, search, minCost, maxCost, minTime, maxTime, sortBy],
    queryFn: async () => {
      const params = new URLSearchParams({
        offset: String((page - 1) * pageSize),
        limit: String(pageSize),
        search,
        sortBy,
      })

      // Добавляем параметры фильтрации по стоимости
      if (minCost !== null) {
        params.append('minCost', String(minCost))
      }
      if (maxCost !== null) {
        params.append('maxCost', String(maxCost))
      }

      // Добавляем параметры фильтрации по времени
      if (minTime !== null) {
        params.append('minTime', String(minTime * 60 * 60 * 1_000_000_000)) // конвертируем часы в наносекунды
      }
      if (maxTime !== null) {
        params.append('maxTime', String(maxTime * 60 * 60 * 1_000_000_000)) // конвертируем часы в наносекунды
      }

      const res = await api.get(`/orders?${params}`)
      return res.data as { orders: Array<Order>; total: number }
    },
  })

  const orders: Array<Order> = data && 'orders' in data ? data.orders : []
  const total = data && 'total' in data ? data.total : 0

  return (
    <div style={{ padding: 24 }}>
      <h1 style={{ fontSize: 24, marginBottom: 16 }}>Главная исполнителя</h1>
      <Row gutter={[16, 16]} align="middle" style={{ marginBottom: 16 }}>
        <Col>
          <Select 
            value={sortBy} 
            onChange={setSortBy} 
            style={{ width: 200 }}
          >
            <Option value="newest">Сначала новые</Option>
            <Option value="oldest">Сначала старые</Option>
            <Option value="cost_asc">По возрастанию цены</Option>
            <Option value="cost_desc">По убыванию цены</Option>
            <Option value="time_asc">По возрастанию времени</Option>
            <Option value="time_desc">По убыванию времени</Option>
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
            <Row gutter={[32, 16]}>
              <Col>
                <Space direction="vertical" size="small">
                  <div>Стоимость:</div>
                  <Space>
                    <Space.Compact>
                      <InputNumber
                        min={0}
                        max={maxCost ?? 1000000}
                        value={minCost}
                        onChange={(value) => {
                          if (value === null || maxCost === null || value <= maxCost) {
                            setMinCost(value)
                          }
                        }}
                        placeholder="От"
                        style={{ width: 120 }}
                        formatter={(value) => value ? `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',') : ''}
                        parser={(value) => value ? Number(value.replace(/\$\s?|(,*)/g, '')) : 0}
                      />
                      {minCost !== null && (
                        <Button
                          icon={<CloseOutlined />}
                          onClick={() => setMinCost(null)}
                          style={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}
                        />
                      )}
                    </Space.Compact>
                    <Space.Compact>
                      <InputNumber
                        min={minCost ?? 0}
                        max={1000000}
                        value={maxCost}
                        onChange={(value) => {
                          if (value === null || minCost === null || value >= minCost) {
                            setMaxCost(value)
                          }
                        }}
                        placeholder="До"
                        style={{ width: 120 }}
                        formatter={(value) => value ? `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',') : ''}
                        parser={(value) => value ? Number(value.replace(/\$\s?|(,*)/g, '')) : 0}
                      />
                      {maxCost !== null && (
                        <Button
                          icon={<CloseOutlined />}
                          onClick={() => setMaxCost(null)}
                          style={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}
                        />
                      )}
                    </Space.Compact>
                  </Space>
                </Space>
              </Col>
              <Col>
                <Space direction="vertical" size="small">
                  <div>Время выполнения:</div>
                  <Space>
                    <Space.Compact>
                      <InputNumber
                        min={1}
                        max={maxTime ?? 720}
                        value={minTime}
                        onChange={(value) => {
                          if (value === null || maxTime === null || value <= maxTime) {
                            setMinTime(value)
                          }
                        }}
                        placeholder="От"
                        style={{ width: 100 }}
                      />
                      {minTime !== null && (
                        <Button
                          icon={<CloseOutlined />}
                          onClick={() => setMinTime(null)}
                          style={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}
                        />
                      )}
                    </Space.Compact>
                    <Space.Compact>
                      <InputNumber
                        min={minTime ?? 1}
                        max={720}
                        value={maxTime}
                        onChange={(value) => {
                          if (value === null || minTime === null || value >= minTime) {
                            setMaxTime(value)
                          }
                        }}
                        placeholder="До"
                        style={{ width: 100 }}
                      />
                      {maxTime !== null && (
                        <Button
                          icon={<CloseOutlined />}
                          onClick={() => setMaxTime(null)}
                          style={{ borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}
                        />
                      )}
                    </Space.Compact>
                  </Space>
                </Space>
              </Col>
            </Row>
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
              <Link to={`/orders/${order.id}` as any}>
                <Card
                  className="order-card"
                  style={{
                    height: '100%',
                    boxShadow: '0 2px 8px rgba(0,0,0,0.09)',
                  }}
                  hoverable
                >
                  <div style={{ fontWeight: 600, marginBottom: 8, fontSize: 16 }}>{order.title}</div>
                  <div style={{
                    marginBottom: 12,
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
                      ⏰ {formatCompletionTime(order.completionTime)}
                    </span>
                    <span style={{
                      fontWeight: 700,
                      fontSize: 18,
                      color: '#1890ff'
                    }}>
                      💰 {order.cost ? `${order.cost.toLocaleString()} ₽` : 'По договорённости'}
                    </span>
                  </div>
                </Card>
              </Link>
            </Col>
          ))}
        </Row>
      )}
      <div style={{ marginTop: 24, textAlign: 'center' }}>
        <Pagination
          current={page}
          pageSize={pageSize}
          total={total}
          onChange={(newPage) => setPage(newPage)}
          onShowSizeChange={(_, newSize) => {
            setPageSize(newSize)
            setPage(1)
          }}
          showSizeChanger
          pageSizeOptions={[6, 12, 18, 24]}
          showTotal={(totalItems) => `Всего ${totalItems} заказов`}
        />
      </div>
    </div>
  )
} 