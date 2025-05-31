import { useState } from 'react'
import { Button, Card, Checkbox, DatePicker, Form, Input, InputNumber, Typography } from 'antd'

const { Title } = Typography

export default function OrderEdit() {
  const [form] = Form.useForm()
  const [isNegotiable, setIsNegotiable] = useState(false)

  const onNegotiableChange = (e: any) => {
    setIsNegotiable(e.target.checked)
    if (e.target.checked) {
      form.setFieldsValue({ cost: null })
    }
  }

  const onFinish = (values: any) => {
    console.log('Submitted order:', values)
  }

  return (
    <div style={{ maxWidth: 600, margin: '40px auto', padding: 24 }}>
      <Title level={3} style={{ marginBottom: 24 }}>Создание заказа</Title>
      <Card>
        <Form
          layout="vertical"
          form={form}
          onFinish={onFinish}
          initialValues={{
            title: '',
            description: '',
            deadline: null,
            cost: null,
            negotiable: false,
          }}
        >
          <Form.Item
            label="Название"
            name="title"
            rules={[{ required: true, message: 'Введите название заказа' }]}
          >
            <Input placeholder="Введите название" />
          </Form.Item>

          <Form.Item
            label="Описание"
            name="description"
            rules={[{ required: true, message: 'Введите описание заказа' }]}
          >
            <Input.TextArea rows={4} placeholder="Введите описание" />
          </Form.Item>

          <Form.Item
            label="Срок выполнения"
            name="deadline"
            rules={[{ required: true, message: 'Выберите срок выполнения' }]}
          >
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item label="Стоимость" style={{ marginBottom: 0 }}>
            <Form.Item
              name="cost"
              rules={[
                {
                  required: !isNegotiable,
                  message: 'Введите стоимость или отметьте договорную',
                },
                {
                  type: 'number',
                  min: 0,
                  message: 'Стоимость не может быть меньше 0',
                },
              ]}
              style={{ display: 'inline-block', width: 'calc(70% - 8px)' }}
            >
              <InputNumber
                min={0}
                disabled={isNegotiable}
                style={{ width: '100%' }}
                placeholder="Введите стоимость"
                formatter={value => `${value} ₽`}
              />
            </Form.Item>

            <Form.Item
              name="negotiable"
              valuePropName="checked"
              style={{ display: 'inline-block', width: '30%', paddingLeft: 8 }}
            >
              <Checkbox checked={isNegotiable} onChange={onNegotiableChange}>
                Договорная
              </Checkbox>
            </Form.Item>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Разместить
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}
