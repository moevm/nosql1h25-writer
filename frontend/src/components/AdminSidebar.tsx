import { Layout, Menu } from 'antd';
import { BarChartOutlined, ExportOutlined, ImportOutlined, UserOutlined } from '@ant-design/icons';
import { useNavigate } from '@tanstack/react-router';
import { useEffect, useState } from 'react';

const { Sider } = Layout;

interface AdminSidebarProps {
  currentPath: string;
}

const AdminSidebar = ({ currentPath }: AdminSidebarProps) => {
  const navigate = useNavigate();
  const [selectedKeys, setSelectedKeys] = useState<Array<string>>([currentPath]);

  useEffect(() => {
    setSelectedKeys([currentPath]);
  }, [currentPath]);

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate({ to: key });
    setSelectedKeys([key]);
  };

  return (
    <Sider
      theme="dark"
      width={250}
      style={{
        overflow: 'auto',
        position: 'fixed',
        left: 0,
        top: 64,
        bottom: 0,
        background: '#001529',
      }}
    >
      <div style={{ padding: '20px', color: 'white' }}>
        <h2 style={{ margin: 0 }}>Админ-панель</h2>
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={selectedKeys}
        onClick={handleMenuClick}
        items={[
          {
            key: '/admin/users',
            icon: <UserOutlined />,
            label: 'Пользователи',
          },
          {
            key: '/admin/import',
            icon: <ImportOutlined />,
            label: 'Импорт базы',
          },
          {
            key: '/admin/export',
            icon: <ExportOutlined />,
            label: 'Экспорт базы',
          },
          {
            key: '/admin/stats',
            icon: <BarChartOutlined />,
            label: 'Статистика',
          },
        ]}
      />
    </Sider>
  );
};

export default AdminSidebar; 