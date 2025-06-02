import { Alert, Button, Card, Upload, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { useState } from 'react';
import type { UploadFile } from 'antd/es/upload/interface';
import type { AxiosError } from 'axios';
import { api } from '@/integrations/api.ts';

interface EchoHTTPError {
  message: string;
}

export const ImportDatabase = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [fileList, setFileList] = useState<Array<UploadFile>>([]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleImport = async () => {
    if (fileList.length === 0) return;

    try {
      setIsLoading(true);
      setError(null);
      setSuccess(false);
      const formData = new FormData();
      formData.append('file', fileList[0].originFileObj as File);

      await api.post('/admin/import', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      setSuccess(true);
      message.success({
        content: 'База данных успешно импортирована',
        duration: 5,
        style: {
          marginTop: '20vh',
        },
      });
      setFileList([]);
    } catch (err) {
      console.error('Error importing database:', err);
      const axiosError = err as AxiosError<EchoHTTPError>;
      console.log('Error response:', {
        status: axiosError.response?.status,
        statusText: axiosError.response?.statusText,
        data: axiosError.response?.data,
        headers: axiosError.response?.headers
      });
      const errorMessage = axiosError.response?.data.message || 'Ошибка при импорте базы данных';
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
    <Card title="Импорт базы данных" style={{ margin: '20px' }}>
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
        {success && (
          <Alert
            message="Успех"
            description="База данных была успешно импортирована"
            type="success"
            showIcon
            closable
            onClose={() => setSuccess(false)}
          />
        )}
        <Upload
          accept=".gzip"
          maxCount={1}
          fileList={fileList}
          onChange={({ fileList: newFileList }) => setFileList(newFileList)}
          beforeUpload={() => false}
        >
          <Button icon={<UploadOutlined />}>Выбрать файл</Button>
        </Upload>
        <div style={{ color: '#666', fontSize: '14px' }}>
          Поддерживаемый формат: GZIP
        </div>
        <Button
          type="primary"
          onClick={handleImport}
          loading={isLoading}
          disabled={fileList.length === 0}
        >
          Импортировать базу данных
        </Button>
      </div>
    </Card>
  );
}; 