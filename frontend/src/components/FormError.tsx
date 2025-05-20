import { Alert } from 'antd'

interface FormErrorProps {
  message: string
}

export function FormError({ message }: FormErrorProps) {
  if (!message) return null
  
  return (
    <Alert
      message={message}
      type="error"
      showIcon
      className="mb-4"
    />
  )
} 