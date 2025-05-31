import { createRoute } from '@tanstack/react-router'
import type { RootRoute } from '@tanstack/react-router'
import ProtectedRoute from '../components/ProtectedRoute'
import OrderEdit from '../components/OrderEdit'

export default (parentRoute: RootRoute) =>
  createRoute({
    path: '/edit-order',
    component: () => (
      <ProtectedRoute allowedRoles={['client']}>
        <OrderEdit />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  })