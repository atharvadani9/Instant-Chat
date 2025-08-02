// noinspection JSIgnoredPromiseFromCall

import type {User} from "../types";
import {useAuthContext} from "../contexts/AuthContext.tsx";
import {useCallback, useEffect, useState} from "react";
import {getAPI} from "../utils/api.ts";

export const useUsers = () => {
    const {userID} = useAuthContext();
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchUsers = useCallback(async () => {
        if (!userID) {
            return;
        }
        setLoading(true);
        setError(null);

        try {
            const resp = await getAPI(`/user.get?user_id=${userID}`);
            if (resp.error && resp.error !== "") {
                setError(resp.error);
            } else {
                setUsers(resp.users || []);
            }
        } catch (err) {
            console.log("Error fetching users:", err);
            setError("Failed to fetch users");
        } finally {
            setLoading(false);
        }
    }, [userID]);

    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    return {
        users,
        loading,
        error,
        fetchUsers,
    };
}