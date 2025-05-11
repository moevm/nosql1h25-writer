import { createRoute } from '@tanstack/react-router';
import type { RootRoute } from '@tanstack/react-router';
import AuthPage from '../components/AuthPage';
import UnauthRoute from '../components/UnauthRoute';

export default (parentRoute: RootRoute) =>
  createRoute({
    path: '/login',
    component: () => (
      <UnauthRoute>
        <AuthPage />
      </UnauthRoute>
    ),
    getParentRoute: () => parentRoute,
  }); 