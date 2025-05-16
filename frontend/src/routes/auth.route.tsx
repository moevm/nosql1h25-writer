import { createRoute } from '@tanstack/react-router';
import AuthPage from '../components/AuthPage';
import UnauthRoute from '../components/UnauthRoute';
import type { RootRoute } from '@tanstack/react-router';

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