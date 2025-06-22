import { lazyRouteComponent } from '@tanstack/react-router';

export const Route = {
  path: '/admin/stats',
  component: lazyRouteComponent(() => import('../components/admin/AdminStatsPage')),
}; 