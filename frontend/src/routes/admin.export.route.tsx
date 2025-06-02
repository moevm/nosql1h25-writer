import { createRoute } from '@tanstack/react-router';
import { ExportDatabase } from '../components/admin/ExportDatabase';
import type { Route } from '@tanstack/react-router';

export const createAdminExportRoute = (parentRoute: Route) =>
  createRoute({
    getParentRoute: () => parentRoute,
    path: 'export',
    component: () => <ExportDatabase />,
  }); 