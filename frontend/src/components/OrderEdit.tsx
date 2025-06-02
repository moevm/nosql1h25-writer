import { Button, Form, Input, InputNumber, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from '@tanstack/react-router';
import { api } from '../integrations/api';

interface EditOrderForm {
  title?: string;
  description?: string;
  completionTime?: number; // в часах
  cost?: number;
}

export default function EditOrderPage() {
  const navigate = useNavigate();
  const params = useParams({ strict: false }) as { id: string };
  const [form] = Form.useForm<EditOrderForm>();
  const [loading, setLoading] = useState(true);

  const [originalOrder, setOriginalOrder] = useState<any>(null);

  useEffect(() => {
    const fetchOrder = async () => {
      try {
        const response = await api.get(`/orders/${params.id}`);
        const order = response.data.order;
        setOriginalOrder(order);

        form.setFieldsValue({
          title: order.title,
          description: order.description,
          completionTime: Math.round(order.completionTime / 3600000000000),
          cost: order.cost,
        });

        setLoading(false);
      } catch (error) {
        message.error('Не удалось загрузить заказ');
        setLoading(false);
      }
    };
    fetchOrder();
  }, [params.id, form]);

  const onFinish = async (values: EditOrderForm) => {
    try {
      const payload: Record<string, any> = {};

      if (values.title?.trim() && values.title !== originalOrder.title) {
        payload.title = values.title.trim();
      }
      if (values.description?.trim() && values.description !== originalOrder.description) {
        payload.description = values.description.trim();
      }
      if (
        values.completionTime &&
        values.completionTime * 3600000000000 !== originalOrder.completionTime
      ) {
        payload.completionTime = values.completionTime * 3600000000000;
      }
      if (
        typeof values.cost === 'number' &&
        values.cost !== originalOrder.cost
      ) {
        payload.cost = values.cost;
      }

      if (Object.keys(payload).length === 0) {
        message.warning('Нет изменений для сохранения');
        return;
      }

      await api.patch(`/orders/${params.id}`, payload);

      message.success('Заказ успешно обновлён');
      navigate({ to: `/orders/${params.id}` });
    } catch (error) {
      message.error('Ошибка при обновлении заказа');
    }
  };

  if (loading || !originalOrder) {
    return <p className="text-center mt-10">Загрузка...</p>;
  }

  return (
    <div className="max-w-2xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Редактирование заказа</h1>

      {!loading && originalOrder && (
        <Form
          form={form}
          layout="vertical"
          onFinish={onFinish}
          className="space-y-4"
        >
          <Form.Item
            label="Название"
            name="title"
            rules={[
              { min: 3, max: 256, message: 'Название должно быть от 3 до 256 символов' },
            ]}
          >
            <Input placeholder="Введите новое название (опционально)" />
          </Form.Item>

          <Form.Item
            label="Описание"
            name="description"
            rules={[
              { min: 16, max: 2048, message: 'Описание должно быть от 16 до 2048 символов' },
            ]}
          >
            <Input.TextArea
              placeholder="Введите новое описание (опционально)"
              rows={6}
            />
          </Form.Item>

          <Form.Item
            label="Время выполнения (часы)"
            name="completionTime"
            rules={[
              { type: 'number', min: 1, message: 'Минимум 1 час' },
            ]}
          >
            <InputNumber min={1} className="w-full" />
          </Form.Item>

          <Form.Item
            label="Стоимость (опционально)"
            name="cost"
            rules={[
              { type: 'number', min: 0, message: 'Стоимость не может быть отрицательной' },
            ]}
          >
            <InputNumber min={0} className="w-full" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" className="w-full">
              Сохранить изменения
            </Button>
          </Form.Item>
        </Form>
      )}
    </div>
  );
}
