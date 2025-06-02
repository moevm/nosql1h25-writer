import { createRoute } from '@tanstack/react-router';
import { Typography } from 'antd';
import { UsersList } from '../components/admin/UsersList';
import type { Route } from '@tanstack/react-router';
import './AdminUsers.css';

const { Title } = Typography;

export const createAdminUsersRoute = (parentRoute: Route) =>
  createRoute({
    path: '/users',
    component: () => (
      <div className="admin-users">
        <Title level={2}>Управление пользователями</Title>
        <UsersList />
      </div>
    ),
    getParentRoute: () => parentRoute,
  }); 