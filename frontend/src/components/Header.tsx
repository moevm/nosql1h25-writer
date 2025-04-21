import { Link } from '@tanstack/react-router'

export default function Header() {
  return (
    <header className="p-3 flex gap-3 bg-white text-black justify-between">
      <nav className="flex flex-row">
        <div className="px-2 font-bold">
          <Link to="/">Home</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to="/auth/login">Login</Link>
        </div>
      </nav>
    </header>
  )
}
