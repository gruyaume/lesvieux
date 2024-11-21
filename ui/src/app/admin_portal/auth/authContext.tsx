"use client"

import { createContext, useContext, useState, useEffect, Dispatch, SetStateAction } from 'react';
import { User } from '../../types';
import { useCookies } from 'react-cookie';
import { jwtDecode } from 'jwt-decode';
import { useRouter } from 'next/navigation';
import { isLoggedIn } from '../../queries';

type AuthContextType = {
    user: User | null
    firstUserCreated: boolean
    setFirstUserCreated: Dispatch<SetStateAction<boolean>>
}

const AuthContext = createContext<AuthContextType>({ user: null, firstUserCreated: false, setFirstUserCreated: () => { } });

export const AuthProvider = ({ children }: Readonly<{ children: React.ReactNode }>) => {
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const [user, setUser] = useState<User | null>(null);
    const [firstUserCreated, setFirstUserCreated] = useState<boolean>(false);
    const router = useRouter();

    useEffect(() => {
        const token = cookies.user_token;
        if (token) {
            let userObject = jwtDecode(cookies.user_token) as User;
            userObject.authToken = cookies.user_token;

            isLoggedIn(token)
                .then((IsLoggedIn) => {
                    if (IsLoggedIn) {
                        setUser(userObject);
                        setFirstUserCreated(true);
                    } else {
                        setUser(null);
                        router.push('/admin_portal/login');
                    }
                })
                .catch((error) => {
                    console.error('Error verifying token:', error);
                    setUser(null);
                    router.push('/admin_portal/login');
                });
        } else {
            setUser(null);
            router.push('/admin_portal/login');
        }
    }, [cookies.user_token, router]);

    return (
        <AuthContext.Provider value={{ user, firstUserCreated, setFirstUserCreated }}>
            {children}
        </AuthContext.Provider>
    );
};



export const useAuth = () => useContext(AuthContext);