import { useAppForm } from '../hooks/demo.form'
import { z } from 'zod'
import { useAuth } from '../context/AuthContext'
import { message } from 'antd'
import { useNavigate } from '@tanstack/react-router'
import { useFieldContext } from '../hooks/demo.form-context'

const schema = z.object({
  displayName: z.string().min(3, 'Минимум 3 символа'),
  email: z.string().email('Введите корректный email'),
  password: z.string().min(8, 'Минимум 8 символов'),
  confirm: z.string(),
}).refine((data) => data.password === data.confirm, {
  message: 'Пароли не совпадают',
  path: ['confirm'],
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

export default function RegisterPage() {
  const { login } = useAuth()
  const navigate = useNavigate()

  const form = useAppForm({
    defaultValues: { displayName: '', email: '', password: '', confirm: '' },
    validators: { onBlur: schema },
    onSubmit: async ({ value }) => {
      try {
        const res = await fetch('/api/auth/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            displayName: value.displayName,
            email: value.email,
            password: value.password,
          }),
        })
        if (!res.ok) throw new Error('Ошибка регистрации')
        const data = await res.json()
        login(data.accessToken, data.refreshToken)
        message.success('Регистрация успешна!')
        navigate({ to: '/' })
      } catch (e) {
        message.error('Ошибка регистрации')
      }
    },
  })

  return (
    <div style={{ maxWidth: 400, margin: 'auto', padding: 32 }}>
      <h2>Зарегистрироваться</h2>
      <form
        onSubmit={e => {
          e.preventDefault()
          e.stopPropagation()
          form.handleSubmit()
        }}
      >
        <form.AppField name="displayName">
          {field => <field.TextField label="Имя" />}
        </form.AppField>
        <form.AppField name="email">
          {field => <field.TextField label="Email" />}
        </form.AppField>
        <form.AppField name="password">
          {() => <PasswordField label="Пароль" />}
        </form.AppField>
        <form.AppField name="confirm">
          {() => <PasswordField label="Повторите пароль" />}
        </form.AppField>
        <form.AppForm>
          <form.SubscribeButton label="Зарегистрироваться" />
        </form.AppForm>
      </form>
    </div>
  )
} 