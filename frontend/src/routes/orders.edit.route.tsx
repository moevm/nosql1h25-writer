import { createRoute } from '@tanstack/react-router'
import ProtectedRoute from '../components/ProtectedRoute'
import EditOrderPage from '../components/OrderEdit'
import type { RootRoute } from '@tanstack/react-router'

export const createOrderEditRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/orders/$id/edit',
    component: () => (
      <ProtectedRoute allowedRoles={['client']}>
        <EditOrderPage />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  })