import { createRoute } from '@tanstack/react-router'
import ProtectedRoute from '../components/ProtectedRoute'
import OrderDetails from '../components/OrderDetails'
import type { RootRoute } from '@tanstack/react-router'

export const createOrderDetailsRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/orders/$id',
    component: () => (
      <ProtectedRoute>
        <OrderDetails />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  })