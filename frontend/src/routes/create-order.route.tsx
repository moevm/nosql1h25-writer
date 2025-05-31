import { createRoute } from '@tanstack/react-router';
import ProtectedRoute from '../components/ProtectedRoute';
import CreateOrderPage from '../components/CreateOrderPage';
import type { RootRoute } from '@tanstack/react-router';

export const createCreateOrderRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/orders/create',
    component: () => (
      <ProtectedRoute allowedRoles={['client']}>
        <CreateOrderPage />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  }); 