import { Navigate, createRoute, useParams } from '@tanstack/react-router';
import UserProfilePage from '../components/UserProfilePage';
import ProtectedRoute from '../components/ProtectedRoute';
import { getUserIdFromToken } from '../integrations/auth';
import type { RootRoute } from '@tanstack/react-router';

export const createUserProfileRoute = (parentRoute: RootRoute) =>
  createRoute({
    getParentRoute: () => parentRoute,
    path: '/users/$userId',
    component: () => {
      const { userId } = useParams({ from: '/users/$userId' });
      const currentUserId = getUserIdFromToken();

      if (userId === currentUserId) {
        return <Navigate to="/profile" />;
      }

      return (
        <ProtectedRoute allowedRoles={['client', 'freelancer']}>
          <UserProfilePage userId={userId} />
        </ProtectedRoute>
      );
    },
  }); 