import { createRoute } from '@tanstack/react-router'
import ProtectedRoute from '../components/ProtectedRoute'
import { OrdersPage } from '../components/OrdersPage'
import type { RootRoute } from '@tanstack/react-router'

export const createOrdersRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/orders',
    component: () => (
      <ProtectedRoute allowedRoles={['freelancer']}>
        <OrdersPage />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  })