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
import { Typography } from 'antd'
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
import { createUserResponsesRoute } from './routes/user-responses.route'
import { createEditProfileRoute } from './routes/edit-profile.route.tsx'
import AdminLayout from './routes/AdminLayout'
import ProfilePage from './components/ProfilePage'
import ProtectedRoute from './components/ProtectedRoute'

import Header from './components/Header'

import TanstackQueryLayout from './integrations/tanstack-query/layout'

import * as TanstackQuery from './integrations/tanstack-query/root-provider'

import 'antd/dist/reset.css'
import './styles.css'
import './styles/global.css'
import reportWebVitals from './reportWebVitals.ts'

import { AuthProvider } from './context/AuthContext'
import { createUserProfileRoute } from './routes/user-profile.route'
import { UsersList } from '@/components/admin/UsersList'
import { ImportDatabase } from '@/components/admin/ImportDatabase'
import { ExportDatabase } from '@/components/admin/ExportDatabase'

const { Title } = Typography;

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
  component: () => <Navigate to="/orders" search={{}} />,
})

const adminRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'admin',
  component: AdminLayout,
})

const adminUsersRoute = createRoute({
  getParentRoute: () => adminRoute,
  path: 'users',
  component: () => (
    <div className="admin-users">
      <Title level={2}>Управление пользователями</Title>
      <UsersList />
    </div>
  ),
})

const adminImportRoute = createRoute({
  getParentRoute: () => adminRoute,
  path: 'import',
  component: () => <ImportDatabase />,
})

const adminExportRoute = createRoute({
  getParentRoute: () => adminRoute,
  path: 'export',
  component: () => <ExportDatabase />,
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
  createUserResponsesRoute(rootRoute),
  createEditProfileRoute(rootRoute),
  createUserProfileRoute(rootRoute),
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
