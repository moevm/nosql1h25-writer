import { createRoute } from '@tanstack/react-router';
import ProtectedRoute from '../components/ProtectedRoute';
import { UserOrders } from '../pages/UserOrders';
import type { RootRoute } from '@tanstack/react-router';

export const createUserOrdersRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/profile/orders',
    component: () => (
      <ProtectedRoute>
        <UserOrders />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  }); 