import { useQuery } from '@tanstack/react-query';
import { api, getUserIdFromToken } from '../integrations/auth';

interface UserProfile {
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