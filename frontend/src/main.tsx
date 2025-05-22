import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import {
  Outlet,
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter,
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import DemoFormAddress from './routes/demo.form.address'
import DemoFormSimple from './routes/demo.form.simple'
import DemoTable from './routes/demo.table'
import DemoTanstackQuery from './routes/demo.tanstack-query'
import AuthRoute from './routes/auth.route'
import RegisterRoute from './routes/register.route'
import OrdersRoute from './routes/orders.route'
import OrderDetailsRoute from './routes/order.details.route'
import ProfileRoute from './routes/profile.route'

import Header from './components/Header'

import TanstackQueryLayout from './integrations/tanstack-query/layout'

import * as TanstackQuery from './integrations/tanstack-query/root-provider'

import './styles.css'
import reportWebVitals from './reportWebVitals.ts'

import App from './App.tsx'
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
  component: App,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  DemoFormAddress(rootRoute),
  DemoFormSimple(rootRoute),
  DemoTable(rootRoute),
  DemoTanstackQuery(rootRoute),
  AuthRoute(rootRoute),
  RegisterRoute(rootRoute),
  OrdersRoute(rootRoute),
  OrderDetailsRoute(rootRoute),
  ProfileRoute(rootRoute),
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
