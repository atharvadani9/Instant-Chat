import {useAuthContext} from "../../contexts/AuthContext.tsx";
import {useWebSocket} from "../../hooks/useWebSocket.ts";
import {useUsers} from "../../hooks/useUsers.ts";
import type {User} from "../../types";
import {useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";
import {Box, Grid, Paper, Typography} from "@mui/material";
import UserList from "./UserList.tsx";
import MessageWindow from "./MessageWindow.tsx";
import UserAvatar from "./UserAvatar.tsx";

export const ChatPage = () => {
    const {userID} = useAuthContext();
    const navigate = useNavigate();
    const {messages, sendMessage, getMessages, clearMessages, connectionState, error} = useWebSocket();
    const {users, loading: usersLoading, error: usersError} = useUsers();

    const [selectedUser, setSelectedUser] = useState<User | null>(null);

    // Handle user authentication
    useEffect(() => {
        if (!userID) {
            navigate("/login");
            return;
        }
    }, [userID, navigate]);

    // Handle initial user selection when users are loaded
    useEffect(() => {
        if (users.length > 0 && !selectedUser) {
            const firstUser = users[0];
            setSelectedUser(firstUser);
            clearMessages();
            getMessages(firstUser.id);
        }
    }, [users, selectedUser]); // Removed getMessages and clearMessages from dependencies

    // Handle getting messages when user is selected and WebSocket is connected
    useEffect(() => {
        if (selectedUser && connectionState === "connected") {
            getMessages(selectedUser.id);
        }
    }, [selectedUser, connectionState]); // Removed getMessages from dependencies

    const handleUserSelect = (user: User) => {
        setSelectedUser(user);
        clearMessages();
        // getMessages will be called by the useEffect above when selectedUser changes
    }

    const handleSendMessage = (content: string) => {
        if (!selectedUser) {
            return;
        }
        sendMessage(selectedUser.id, content);
        // Don't call getMessages here - real-time messages will come through WebSocket
    }

    return (
        <Grid
            display="flex"
            flexDirection="column"
            justifyContent="flex-start"
            sx={{minHeight: '100vh', width: '100%'}}
        >
            <Paper elevation={1} sx={{p: 2, borderRadius: 0, borderColor: 'divider', borderBottom: 1, width: '100%'}}>
                <Box sx={{display: 'flex', justifyContent: 'center', alignItems: 'center', position: 'relative'}}>
                    <Typography variant="h5" component="h1">{"Instant Chat"}</Typography>
                    {selectedUser && (
                        <Box sx={{position: 'absolute', right: 0}}>
                            <UserAvatar/>
                        </Box>
                    )}
                </Box>
            </Paper>
            <Grid container sx={{flexGrow: 1, overflow: 'hidden', p: 1, display: 'flex', flexDirection: 'row',}}>
                <Grid size={3} sx={{borderColor: 'divider', borderRight: 1, p: 1}}>
                    <UserList
                        users={users}
                        onUserSelect={handleUserSelect}
                        selectedUser={selectedUser}
                        error={usersError}
                        loading={usersLoading}
                    />
                </Grid>
                <Grid size={9}>
                    <MessageWindow
                        messages={messages}
                        onSendMessage={handleSendMessage}
                        error={error}
                        connectionState={connectionState}
                        selectedUser={selectedUser}
                    />
                </Grid>
            </Grid>
        </Grid>
    )
}

export default ChatPage;
