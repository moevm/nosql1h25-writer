import {createFileRoute, Navigate, useNavigate} from '@tanstack/react-router'
import {type FormEvent, useState} from "react";
import {useAuth} from "@/auth/AuthProvider.tsx";

export const Route = createFileRoute('/login')({
    component: Login,
})

function Login() {
    const { login, user } = useAuth();
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleLogin = async (e: FormEvent)=> {
        e.preventDefault();
        try {
            await login(email, password);
            await navigate({to: '/orders'});
        } catch (error) {
            alert(error.message);
        }
    };

    if (user) {
        return <Navigate to="/orders" />;
    }

    return (
        <form onSubmit={handleLogin}>
            <input value={email} onChange={e => setEmail(e.target.value)}/>
            <input type="password" value={password} onChange={e => setPassword(e.target.value)}/>
            <button type="submit">Войти</button>
        </form>
    )
}