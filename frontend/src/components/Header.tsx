import { Link } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { isAuthenticated } from '../integrations/auth'
import LogoutButton from './LogoutButton'

export default function Header() {
  const [auth, setAuth] = useState(isAuthenticated())

  useEffect(() => {
    const check = () => setAuth(isAuthenticated())
    window.addEventListener('auth-changed', check)
    return () => window.removeEventListener('auth-changed', check)
  }, [])

  useEffect(() => {
    setAuth(isAuthenticated())
  })

  return (
    <header className="p-2 flex gap-2 bg-white text-black justify-between items-center">
      <nav className="flex flex-row">
        <div className="px-2 font-bold">
          <Link to="/">Home</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to="/demo/form/simple">Simple Form</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to="/demo/form/address">Address Form</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to="/demo/table">TanStack Table</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to="/demo/tanstack-query">TanStack Query</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to="/orders">Заказы</Link>
        </div>
      </nav>
      {auth && <LogoutButton />}
    </header>
  )
}
