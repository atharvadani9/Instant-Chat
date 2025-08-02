import { useState, useEffect } from 'react';
import { getUserID } from '../utils/utils';

export interface AuthState {
    userID: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
}

export const useAuth = () => {
    const [authState, setAuthState] = useState<AuthState>({
        userID: null,
        isAuthenticated: false,
        isLoading: true,
    });

    useEffect(() => {
        // Automatically get userID on component mount
        const initializeAuth = () => {
            const userID = getUserID();
            setAuthState({
                userID,
                isAuthenticated: !!userID,
                isLoading: false,
            });
        };

        initializeAuth();
    }, []);

    const setUserID = (userID: string) => {
        localStorage.setItem('userID', userID);
        setAuthState({
            userID,
            isAuthenticated: true,
            isLoading: false,
        });
    };

    const logout = () => {
        localStorage.removeItem('userID');
        setAuthState({
            userID: null,
            isAuthenticated: false,
            isLoading: false,
        });
    };

    return {
        ...authState,
        setUserID,
        logout,
    };
};
