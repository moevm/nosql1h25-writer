import { createContext, useContext, useEffect, useState} from 'react';
import { authApi, setAccessToken } from '../api/authApi';
import type {ReactNode} from 'react';

type User = {
    id: string;
    email: string;
};

type AuthContextType = {
    user: User | null;
    loading: boolean;
    login: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | null>(null);

type Props = {
    children: ReactNode;
};

export const AuthProvider = ({ children }: Props) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    const fetchUser = async () => {
        try {
            const res = await authApi.getProfile();
            setUser(res.data);
        } catch {
            setUser(null);
        } finally {
            setLoading(false);
        }
    };

    const login = async (email: string, password: string) => {
        const res = await authApi.login({ email, password });
        setAccessToken(res.data.accessToken);
        await fetchUser();
    };

    const logout = async () => {
        await authApi.logout();
        setAccessToken('');
        setUser(null);
    };

    useEffect(() => {
        (async () => {
            await fetchUser();
        })();
    }, []);

    return (
        <AuthContext.Provider value={{ user, login, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
      throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};