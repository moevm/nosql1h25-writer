import { createRoute } from '@tanstack/react-router';
import ProtectedRoute from '../components/ProtectedRoute';
import { UserResponses } from '../pages/UserResponses';
import type { RootRoute } from '@tanstack/react-router';

export const createUserResponsesRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/profile/responses',
    component: () => (
      <ProtectedRoute allowedRoles={['freelancer']}>
        <UserResponses />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  }); 