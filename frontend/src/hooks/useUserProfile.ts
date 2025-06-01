import { useQuery } from '@tanstack/react-query';
import { getUserIdFromToken } from '../integrations/auth';
import { api } from '../integrations/api';

interface UserProfile {
  id: string;
  displayName: string;
  email: string;
  balance: number;
  client?: {
    rating: number;
    completedOrders: number;
  };
  freelancer?: {
    rating: number;
    completedOrders: number;
  };
}

export function useUserProfile() {
  const userId = getUserIdFromToken();

  return useQuery<UserProfile>({
    queryKey: ['userProfile', userId],
    queryFn: async () => {
      const response = await api.get(`/users/${userId}`);
      return response.data;
    },
    enabled: !!userId,
  });
} 