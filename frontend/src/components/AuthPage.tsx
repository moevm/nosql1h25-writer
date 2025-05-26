import { z } from 'zod'
import { message } from 'antd'
import { Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { useAppForm } from '../hooks/demo.form'
import { useFieldContext } from '../hooks/demo.form-context'
import { FormError } from './FormError'
import { RoleSelector } from './RoleSelector'

const schema = z.object({
  email: z.string().email('Введите корректный email'),
  password: z.string().min(8, 'Минимум 8 символов'),
})

function TextField({ label, type = 'text', hint }: { label: string, type?: string, hint?: string }) {
  const field = useFieldContext<string>()
  return (
    <div className="mb-4">
      <label htmlFor={label} className="block font-medium text-gray-700 mb-1">
        {label}
      </label>
      <input
        type={type}
        value={field.state.value}
        onBlur={field.handleBlur}
        onChange={(e) => field.handleChange(e.target.value)}
        className="w-full px-3 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
      />
      {hint && <div className="text-sm text-gray-500 mt-1">{hint}</div>}
    </div>
  )
}

export default function AuthPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [serverError, setServerError] = useState('')

  const form = useAppForm({
    defaultValues: { 
      email: '', 
      password: '',
    },
    validators: { onBlur: schema },
    onSubmit: async ({ value }) => {
      try {
        setServerError('')
        const res = await fetch('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: value.email,
            password: value.password,
          }),
        })
        
        const data = await res.json()
        
        if (!res.ok) {
          setServerError(data.message || 'Ошибка авторизации')
          return
        }
        
        login(data.accessToken, data.refreshToken)
        message.success('Вход выполнен успешно!')
        await navigate({ to: '/profile' })
      } catch (e) {
        setServerError(e instanceof Error ? e.message : 'Ошибка авторизации')
      }
    },
  })

  return (
    <div className="min-h-[calc(100vh-64px)] flex items-center justify-center bg-gray-50 py-6 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-6 bg-white p-6 rounded-xl shadow-lg">
        <div>
          <h2 className="text-center text-2xl font-bold text-gray-900">
            Войти
          </h2>
        </div>
        <FormError message={serverError} />
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault()
            e.stopPropagation()
            form.handleSubmit()
          }}
        >
          <div className="space-y-4">
            <form.AppField name="email">
              {() => <TextField label="Email" />}
            </form.AppField>
            <form.AppField name="password">
              {() => <TextField label="Пароль" type="password" hint="Минимум 8 символов" />}
            </form.AppField>
          </div>

          <RoleSelector />

          <div className="space-y-4">
            <div className="space-y-1">
              <Link
                to="/register"
                className="block text-blue-600 hover:text-blue-800 transition-colors"
              >
                <span className="border-b border-dotted border-blue-600 hover:border-blue-800">
                  Создать аккаунт
                </span>
              </Link>
            </div>
            <form.AppForm>
              <button
                type="submit"
                className="w-full px-6 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors disabled:opacity-50 text-center"
              >
                Войти
              </button>
            </form.AppForm>
          </div>
        </form>
      </div>
    </div>
  )
} 