import { Button, Card, Input, Space, Table } from 'antd';
import { useState } from 'react';
import type { ColumnsType } from 'antd/es/table';

const { Search } = Input;

interface Order {
  id: string;
  title: string;
  description: string;
  cost: number;
  completionTime: number;
  rating: number;
  clientName: string;
}

const columns: ColumnsType<Order> = [
  {
    title: 'Название',
    dataIndex: 'title',
    key: 'title',
  },
  {
    title: 'Описание',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: 'Стоимость',
    dataIndex: 'cost',
    key: 'cost',
    render: (cost: number) => `${cost.toLocaleString()} руб.`,
  },
  {
    title: 'Срок выполнения',
    dataIndex: 'completionTime',
    key: 'completionTime',
    render: (time: number) => `${time} дней`,
  },
  {
    title: 'Рейтинг',
    dataIndex: 'rating',
    key: 'rating',
    render: (rating: number) => `${rating.toFixed(1)} ⭐`,
  },
  {
    title: 'Заказчик',
    dataIndex: 'clientName',
    key: 'clientName',
  },
];

export default function OrdersList() {
  const [searchText, setSearchText] = useState('');

  // Здесь будет загрузка данных с сервера
  const data: Array<Order> = [];

  return (
    <div style={{ padding: '24px' }}>
      <Card style={{ marginBottom: '24px' }}>
        <Space style={{ width: '100%', justifyContent: 'space-between' }}>
          <Search
            placeholder="Поиск заказов"
            allowClear
            enterButton="Поиск"
            size="large"
            style={{ width: '400px' }}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <Button type="primary" size="large">
            Создать заказ
          </Button>
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showTotal: (total) => `Всего ${total} заказов`,
        }}
      />
    </div>
  );
} 