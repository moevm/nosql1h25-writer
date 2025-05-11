import { Button } from 'antd';
import { useNavigate } from '@tanstack/react-router';
import { api, clearTokens } from '../integrations/auth';

export default function LogoutButton() {
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await api.post('/auth/logout', { refreshToken: localStorage.getItem('refreshToken') });
    } catch {}
    clearTokens();
    navigate({ to: '/login' });
  };

  return (
    <Button onClick={handleLogout} danger>
      Выйти
    </Button>
  );
} 