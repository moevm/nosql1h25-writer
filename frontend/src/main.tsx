import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import {
  Navigate,
  Outlet,
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import DemoFormAddress from './routes/demo.form.address'
import DemoFormSimple from './routes/demo.form.simple'
import DemoTable from './routes/demo.table'
import DemoTanstackQuery from './routes/demo.tanstack-query'
import AuthRoute from './routes/auth.route'
import RegisterRoute from './routes/register.route'
import { createOrdersRoute } from './routes/orders.route'
import { createOrderDetailsRoute } from './routes/order.details.route'
import { createCreateOrderRoute } from './routes/create-order.route'
import { createOrderEditRoute } from './routes/orders.edit.route.tsx'
import { createUserOrdersRoute } from './routes/user-orders.route'
import AdminLayout from './routes/AdminLayout'
import AdminUsers from './routes/AdminUsers'
import { AdminImportRoute } from './routes/admin.import.route'
import { AdminExportRoute } from './routes/admin.export.route'
import ProfilePage from './components/ProfilePage'
import ProtectedRoute from './components/ProtectedRoute'

import Header from './components/Header'

import TanstackQueryLayout from './integrations/tanstack-query/layout'

import * as TanstackQuery from './integrations/tanstack-query/root-provider'

import 'antd/dist/reset.css'
import './styles.css'
import reportWebVitals from './reportWebVitals.ts'

import { AuthProvider } from './context/AuthContext'

const rootRoute = createRootRoute({
  component: () => (
    <>
      <Header />
      <Outlet />
      <TanStackRouterDevtools />

      <TanstackQueryLayout />
    </>
  ),
})

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: () => <Navigate to="/orders" />,
})

const adminRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin',
  component: AdminLayout,
})

const adminUsersRoute = createRoute({
  getParentRoute: () => adminRoute,
  path: '/users',
  component: AdminUsers,
})

const adminImportRoute = createRoute({
  getParentRoute: () => adminRoute,
  path: '/import',
  component: AdminImportRoute,
})

const adminExportRoute = createRoute({
  getParentRoute: () => adminRoute,
  path: '/export',
  component: AdminExportRoute,
})

const profileRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/profile',
  component: () => (
    <ProtectedRoute allowedRoles={['client', 'freelancer']}>
      <ProfilePage />
    </ProtectedRoute>
  ),
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  adminRoute.addChildren([
    adminUsersRoute,
    adminImportRoute,
    adminExportRoute,
  ]),
  DemoFormAddress(rootRoute),
  DemoFormSimple(rootRoute),
  DemoTable(rootRoute),
  DemoTanstackQuery(rootRoute),
  AuthRoute(rootRoute),
  RegisterRoute(rootRoute),
  createOrdersRoute(rootRoute),
  createOrderDetailsRoute(rootRoute),
  createCreateOrderRoute(rootRoute),
  createOrderEditRoute(rootRoute),
  createUserOrdersRoute(rootRoute),
  profileRoute,
])

const router = createRouter({
  routeTree,
  context: {
    ...TanstackQuery.getContext(),
  },
  defaultPreload: 'intent',
  scrollRestoration: true,
  defaultStructuralSharing: true,
  defaultPreloadStaleTime: 0,
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const rootElement = document.getElementById('app')
if (rootElement && !rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <AuthProvider>
        <TanstackQuery.Provider>
          <RouterProvider router={router} />
        </TanstackQuery.Provider>
      </AuthProvider>
    </StrictMode>,
  )
}

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
