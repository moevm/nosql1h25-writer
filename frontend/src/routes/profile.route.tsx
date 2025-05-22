import { createRoute } from '@tanstack/react-router';
import ProfilePage from '../components/ProfilePage';
import ProtectedRoute from '../components/ProtectedRoute';
import type { RootRoute } from '@tanstack/react-router';

function createProfileRoute(parentRoute: RootRoute) {
  return createRoute({
    path: '/profile',
    component: () => (
      <ProtectedRoute allowedRoles={['client', 'freelancer']}>
        <ProfilePage />
      </ProtectedRoute>
    ),
    getParentRoute: () => parentRoute,
  });
}

export default createProfileRoute; 