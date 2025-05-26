import { Layout, Spin } from 'antd';
import { Outlet, useNavigate, useRouter } from '@tanstack/react-router';
import { useEffect } from 'react';
import AdminSidebar from '../components/AdminSidebar';
import { useAuth } from '../context/AuthContext';

const { Content } = Layout;

const AdminLayout = () => {
  const { user, loading } = useAuth();
  const navigate = useNavigate();
  const router = useRouter();
  const currentPath = router.state.location.pathname;

  useEffect(() => {
    console.log('AdminLayout rendering. Current path:', currentPath);
  }, [currentPath]);

  // Пока данные пользователя загружаются, показываем спиннер
  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <Spin size="large" />
      </div>
    );
  }

  // После загрузки проверяем права администратора
  if (!user || user.systemRole !== 'admin') {
    navigate({ to: '/', replace: true });
    return null;
  }

  return (
    <Layout style={{ minHeight: '100vh', marginTop: 64 }}>
      <AdminSidebar currentPath={currentPath} />
      {/* Контентная область справа от сайдбара */}
      <Layout style={{ marginLeft: 250 }}>
        <Content style={{ margin: '0 16px 24px 16px', padding: 24, background: '#fff', minHeight: 'calc(100vh - 64px - 24px - 24px)' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default AdminLayout; 