import { useQuery } from '@tanstack/react-query';
import { api } from '../integrations/api';

export interface UserProfile {
  id: string;
  displayName: string;
  email: string;
  balance: number;
  client?: {
    description: string;
    rating: number;
    updatedAt: string;
  };
  freelancer?: {
    description: string;
    rating: number;
    updatedAt: string;
  };
}

export function useUserProfile(userId: string, role?: 'client' | 'freelancer') {
  return useQuery<UserProfile>({
    queryKey: ['userProfile', userId, role],
    queryFn: async () => {
      const response = await api.get(`/users/${userId}${role ? `?profile=${role}` : ''}`);
      return response.data;
    },
    enabled: !!userId
  });
}
