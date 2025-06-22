import { createRoute } from '@tanstack/react-router';
import ProtectedRoute from '../components/ProtectedRoute';
import type { RootRoute } from '@tanstack/react-router';
import EditProfilePage from '@/components/EditProfilePage';

export const createEditProfileRoute = (parentRoute: RootRoute) =>
  createRoute({
    path: '/profile/edit',
    component: () => (
      <ProtectedRoute allowedRoles={['client', 'freelancer']}>
        <EditProfilePage />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  });