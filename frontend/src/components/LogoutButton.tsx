import { Button } from 'antd';
import { useNavigate } from '@tanstack/react-router';
import { clearTokens } from '../integrations/auth';
import { api } from '../integrations/api';

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