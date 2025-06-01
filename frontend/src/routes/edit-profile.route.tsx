import { createRoute } from '@tanstack/react-router';
import ProtectedRoute from '../components/ProtectedRoute';
import EditProfilePage from '@/components/EditProfilePage';
import type { RootRoute } from '@tanstack/react-router';

export const createEditProfileRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/user/edit',
    component: () => (
      <ProtectedRoute allowedRoles={['client', 'freelancer']}>
        <EditProfilePage />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  });