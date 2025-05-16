import { createRoute } from '@tanstack/react-router';
import RegisterPage from '../components/RegisterPage';
import UnauthRoute from '../components/UnauthRoute';
import type { RootRoute } from '@tanstack/react-router';

export default (parentRoute: RootRoute) =>
  createRoute({
    path: '/register',
    component: () => (
      <UnauthRoute>
        <RegisterPage />
      </UnauthRoute>
    ),
    getParentRoute: () => parentRoute,
  }); 