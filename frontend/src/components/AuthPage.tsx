import { z } from 'zod'
import { useAuth } from '../context/AuthContext'
import { message } from 'antd'
import { useNavigate } from '@tanstack/react-router'
import { useAppForm } from '../hooks/demo.form'
import { useFieldContext } from '../hooks/demo.form-context'

const schema = z.object({
  email: z.string().email('Введите корректный email'),
  password: z.string().min(8, 'Минимум 8 символов'),
})

function PasswordField({ label }: { label: string }) {
  const field = useFieldContext<string>()
  return (
    <div>
      <label htmlFor={label} className="block font-bold mb-1 text-xl">
        {label}
        <input
          type="password"
          value={field.state.value}
          onBlur={field.handleBlur}
          onChange={(e) => field.handleChange(e.target.value)}
          className="w-full px-4 py-2 rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500"
        />
      </label>
    </div>
  )
}

export default function AuthPage() {
  const { login } = useAuth()
  const navigate = useNavigate()

  const form = useAppForm({
    defaultValues: { email: '', password: '' },
    validators: { onBlur: schema },
    onSubmit: async ({ value }) => {
      try {
        const res = await fetch('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(value),
        })
        if (!res.ok) throw new Error('Ошибка авторизации')
        const data = await res.json()
        login(data.accessToken, data.refreshToken)
        message.success('Успешный вход!')
        navigate({ to: '/' })
      } catch (e) {
        message.error('Ошибка авторизации')
      }
    },
  })

  return (
    <div style={{ maxWidth: 400, margin: 'auto', padding: 32 }}>
      <h2>Войти</h2>
      <form
        onSubmit={e => {
          e.preventDefault()
          e.stopPropagation()
          form.handleSubmit()
        }}
      >
        <form.AppField name="email">
          {field => <field.TextField label="Email" />}
        </form.AppField>
        <form.AppField name="password">
          {() => <PasswordField label="Пароль" />}
        </form.AppField>
        <form.AppForm>
          <form.SubscribeButton label="Войти" />
        </form.AppForm>
      </form>
    </div>
  )
} 