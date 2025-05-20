import { Radio } from 'antd'
import { useState } from 'react'
import { roleUtils } from '../utils/role'
import type { UserRole } from '../utils/role'

interface RoleSelectorProps {
  onRoleSelect?: (role: UserRole) => void
}

export function RoleSelector({ onRoleSelect }: RoleSelectorProps) {
  const [selectedRole, setSelectedRole] = useState<UserRole>(roleUtils.getRole())

  const handleRoleSelect = (role: UserRole) => {
    setSelectedRole(role)
    roleUtils.setRole(role)
    onRoleSelect?.(role)
  }

  return (
    <div className="mb-4">
      <label className="block font-medium text-gray-700 mb-1">
        Выберите роль
      </label>
      <Radio.Group 
        value={selectedRole} 
        onChange={(e) => handleRoleSelect(e.target.value)}
        className="w-full"
      >
        <Radio.Button value="client" className="w-1/2 text-center">Заказчик</Radio.Button>
        <Radio.Button value="freelancer" className="w-1/2 text-center">Исполнитель</Radio.Button>
      </Radio.Group>
    </div>
  )
} 