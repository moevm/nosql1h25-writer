import { useQuery } from '@tanstack/react-query';
import { getUserIdFromToken } from '../integrations/auth';
import { api } from '../integrations/api';

interface UserProfile {
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

export function useUserProfile() {
  const userId = getUserIdFromToken();

  return useQuery<UserProfile>({
    queryKey: ['userProfile', userId],
    queryFn: async () => {
      const params = new URLSearchParams();
      params.append('profile', 'client');
      params.append('profile', 'freelancer');

      const response = await api.get(`/users/${userId}?${params.toString()}`);
      console.log(`ANSWER: ${response.data.client}`)
      return response.data;
    },
    enabled: !!userId,
  });
}
