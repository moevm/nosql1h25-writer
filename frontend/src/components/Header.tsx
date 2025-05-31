import { Link } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { Button } from 'antd'
import { isAdmin, isAuthenticated } from '../integrations/auth'
import LogoutButton from './LogoutButton'

export default function Header() {
  const [auth, setAuth] = useState(isAuthenticated())
  const [admin, setAdmin] = useState(isAdmin())

  useEffect(() => {
    const check = () => {
      setAuth(isAuthenticated())
      setAdmin(isAdmin())
    }
    window.addEventListener('auth-changed', check)
    return () => window.removeEventListener('auth-changed', check)
  }, [])

  useEffect(() => {
    setAuth(isAuthenticated())
    setAdmin(isAdmin())
  })

  return (
    <header className="p-2 flex gap-2 bg-white text-black justify-between items-center fixed top-0 left-0 right-0 z-10 h-16">
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

        {admin && (
          <div className="px-2 font-bold">
            <Link to="/admin">Админ-панель</Link>
          </div>
        )}
      </nav>
      <div className="flex items-center gap-2">
        {auth && (
          <>
            <Link to="/profile">
              <Button type="primary">Профиль</Button>
            </Link>
            <LogoutButton />
          </>
        )}
      </div>
    </header>
  )
}
