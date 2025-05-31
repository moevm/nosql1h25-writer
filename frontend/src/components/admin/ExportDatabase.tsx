import { Alert, Button, Card, message } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import { useState } from 'react';
import type { AxiosError } from 'axios';
import { api } from '@/integrations/api.ts';

interface EchoHTTPError {
  message: string;
}

export const ExportDatabase = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleExport = async () => {
    try {
      setIsLoading(true);
      setError(null);
      console.log('Начинаем экспорт базы данных...');
      
      const response = await api.get('/admin/export', {
        responseType: 'blob',
      });

      console.log('Получен ответ от сервера:', {
        status: response.status,
        headers: response.headers,
        data: response.data
      });

      if (!response.data || response.data.size === 0) {
        const errorMessage = 'Получен пустой файл от сервера';
        setError(errorMessage);
        message.error({
          content: errorMessage,
          duration: 5,
          style: {
            marginTop: '20vh',
          },
        });
        setIsLoading(false);
        return;
      }

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'database.gzip');
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);

      message.success('База данных успешно экспортирована');
    } catch (apiError) {
      console.error('Ошибка при экспорте базы данных:', apiError);
      
      const axiosError = apiError as AxiosError<EchoHTTPError>;
      console.log('Детали ошибки:', {
        status: axiosError.response?.status,
        statusText: axiosError.response?.statusText,
        data: axiosError.response?.data,
        headers: axiosError.response?.headers,
        config: {
          url: axiosError.config?.url,
          method: axiosError.config?.method,
          headers: axiosError.config?.headers
        }
      });

      let errorMessage = 'Ошибка при экспорте базы данных';
      
      if (axiosError.response) {
        if (axiosError.response.status === 404) {
          errorMessage = 'Сервер не найден. Проверьте подключение к серверу.';
        } else if (axiosError.response.status === 401) {
          errorMessage = 'Требуется авторизация. Пожалуйста, войдите в систему.';
        } else if (axiosError.response.status === 403) {
          errorMessage = 'Нет прав для экспорта базы данных.';
        } else if (axiosError.response.data.message) {
          errorMessage = axiosError.response.data.message;
        }
      } else if (axiosError.request) {
        errorMessage = 'Нет ответа от сервера. Проверьте подключение к интернету.';
      }

      setError(errorMessage);
      message.error({
        content: errorMessage,
        duration: 5,
        style: {
          marginTop: '20vh',
        },
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card title="Экспорт базы данных" style={{ margin: '10px 20px' }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', maxWidth: '400px' }}>
        {error && (
          <Alert
            message="Ошибка"
            description={error}
            type="error"
            showIcon
            closable
            onClose={() => setError(null)}
          />
        )}
        <div style={{ color: '#666', fontSize: '14px' }}>
          База данных будет экспортирована в формате GZIP
        </div>
        <Button
          type="primary"
          icon={<DownloadOutlined />}
          onClick={handleExport}
          loading={isLoading}
        >
          Экспортировать базу данных
        </Button>
      </div>
    </Card>
  );
}; 