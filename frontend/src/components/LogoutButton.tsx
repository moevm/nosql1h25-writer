import { Button } from 'antd';
import { clearTokens, api } from '../integrations/auth';
import { useNavigate } from '@tanstack/react-router';

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