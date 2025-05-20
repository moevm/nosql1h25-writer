export type UserRole = 'client' | 'freelancer'

const ROLE_STORAGE_KEY = 'user_role'

export const roleUtils = {
  getRole: (): UserRole => {
    const storedRole = localStorage.getItem(ROLE_STORAGE_KEY)
    if (storedRole === 'client' || storedRole === 'freelancer') {
      return storedRole
    }
    return 'client'
  },

  setRole: (role: UserRole): void => {
    localStorage.setItem(ROLE_STORAGE_KEY, role)
  },

  isClient: (): boolean => {
    return roleUtils.getRole() === 'client'
  },

  isFreelancer: (): boolean => {
    return roleUtils.getRole() === 'freelancer'
  }
} 