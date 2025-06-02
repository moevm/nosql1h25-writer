import { Button, Form, Input, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from '@tanstack/react-router'
import { useUserProfile } from '../hooks/useUserProfile';
import { api } from '../integrations/api';
import { getUserIdFromToken } from '../integrations/auth';
import {  roleUtils } from '../utils/role';
import type {UserRole} from '../utils/role';

interface EditProfileForm {
  displayName?: string;
  clientDescription?: string;
  freelancerDescription?: string;
}

export default function EditProfilePage() {
  const navigate = useNavigate();
  const userId = getUserIdFromToken();
  const [selectedRole] = useState<UserRole>(roleUtils.getRole());
  const { data: user, isLoading, refetch } = useUserProfile(String(userId), selectedRole);
  const [form] = Form.useForm<EditProfileForm>();

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        displayName: user.displayName,
        clientDescription: user.client?.description,
        freelancerDescription: user.freelancer?.description,
      });
    }
  }, [user, form]);

  const onFinish = async (values: EditProfileForm) => {
    try {
      const payload: EditProfileForm = {};
      if (values.displayName?.trim()) payload.displayName = values.displayName.trim();
      if (selectedRole === 'client' && values.clientDescription?.trim())
        payload.clientDescription = values.clientDescription.trim();
      if (selectedRole === 'freelancer' && values.freelancerDescription?.trim())
        payload.freelancerDescription = values.freelancerDescription.trim();

      if (Object.keys(payload).length === 0) {
        message.warning('Нет изменений для отправки');
        return;
      }

      await api.patch(`/users/${userId}`, payload);

      message.success('Профиль обновлён');
      await refetch();
      navigate({ to: '/profile' });
    } catch (error) {
      message.error('Ошибка при обновлении профиля');
    }
  };

  if (isLoading) return <p className="text-center mt-10">Загрузка профиля...</p>;

  return (
    <div className="max-w-2xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Редактировать профиль</h1>

      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        className="space-y-4"
      >
        <Form.Item
          label="Отображаемое имя"
          name="displayName"
          rules={[
            { min: 3, max: 64, message: 'Имя должно быть от 3 до 64 символов' }
          ]}
        >
          <Input placeholder="Введите новое имя (опционально)" />
        </Form.Item>

        {selectedRole === 'client' && (
          <Form.Item
            label="Описание как клиент"
            name="clientDescription"
            rules={[
              { min: 16, max: 2048, message: 'Описание должно быть от 16 до 2048 символов' }
            ]}
          >
            <Input.TextArea rows={4} placeholder="Описание клиента (опционально)" />
          </Form.Item>
        )}

        {selectedRole === 'freelancer' && (
          <Form.Item
            label="Описание как фрилансер"
            name="freelancerDescription"
            rules={[
              { min: 16, max: 2048, message: 'Описание должно быть от 16 до 2048 символов' }
            ]}
          >
            <Input.TextArea rows={4} placeholder="Описание фрилансера (опционально)" />
          </Form.Item>
        )}

        <Form.Item>
          <Button type="primary" htmlType="submit" className="w-full">
            Сохранить изменения
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}