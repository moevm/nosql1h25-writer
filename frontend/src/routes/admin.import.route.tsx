import { createRoute } from '@tanstack/react-router';
import { ImportDatabase } from '../components/admin/ImportDatabase';
import type { Route } from '@tanstack/react-router';

export const createAdminImportRoute = (parentRoute: Route) =>
  createRoute({
    path: '/import',
    component: () => <ImportDatabase />,
    getParentRoute: () => parentRoute,
  }); 