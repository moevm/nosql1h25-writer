import { Button, Form, Input, InputNumber, message } from 'antd';
import { useNavigate } from '@tanstack/react-router';
import { api } from '../integrations/api';

interface CreateOrderForm {
  title: string;
  description: string;
  completionTime: number;
  cost?: number;
}

export default function CreateOrderPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm<CreateOrderForm>();

  const onFinish = async (values: CreateOrderForm) => {
    try {
      const response = await api.post('/orders', {
        ...values,
        completionTime: values.completionTime * 3600000000000, // Convert hours to nanoseconds
      });
      
      message.success('Заказ успешно создан');
      navigate({ to: `/orders/${response.data.id}` });
    } catch (error) {
      message.error('Ошибка при создании заказа');
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Создание нового заказа</h1>
      
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
            { required: true, message: 'Пожалуйста, введите название заказа' },
            { min: 3, max: 32, message: 'Название должно быть от 3 до 32 символов' }
          ]}
        >
          <Input placeholder="Введите название заказа" />
        </Form.Item>

        <Form.Item
          label="Описание"
          name="description"
          rules={[
            { required: true, message: 'Пожалуйста, введите описание заказа' },
            { min: 16, max: 8192, message: 'Описание должно быть от 16 до 8192 символов' }
          ]}
        >
          <Input.TextArea 
            placeholder="Введите подробное описание заказа"
            rows={6}
          />
        </Form.Item>

        <Form.Item
          label="Время выполнения (часы)"
          name="completionTime"
          rules={[
            { required: true, message: 'Пожалуйста, укажите время выполнения' },
            { type: 'number', min: 1, message: 'Минимальное время выполнения - 1 час' }
          ]}
        >
          <InputNumber min={1} className="w-full" />
        </Form.Item>

        <Form.Item
          label="Стоимость (опционально)"
          name="cost"
          rules={[
            { type: 'number', min: 0, message: 'Стоимость не может быть отрицательной' }
          ]}
        >
          <InputNumber min={0} className="w-full" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" className="w-full">
            Создать заказ
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
} 