import React, {useEffect, useRef, useState} from "react";
import type {ConnectionState, User, WSMessage} from "../../types";
import {Avatar, Box, CircularProgress, Grid, Paper, TextField, Tooltip, Typography} from "@mui/material";
import {enqueueSnackbar} from "notistack";
import IconButton from "@mui/material/IconButton";
import SendIcon from '@mui/icons-material/Send';
import {getAvatarColor} from "../../utils/utils.ts";
import {useAuthContext} from "../../contexts/AuthContext.tsx";

type MessageWindowProps = {
    selectedUser: User | null;
    messages: WSMessage[];
    onSendMessage: (content: string) => void;
    connectionState: ConnectionState;
    error: string | null;
}

const MessageWindow: React.FC<MessageWindowProps> = ({
                                                         selectedUser,
                                                         messages,
                                                         onSendMessage,
                                                         connectionState,
                                                         error
                                                     }) => {
    const {userID} = useAuthContext();
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const [inputText, setInputText] = useState("");

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({behavior: "smooth"});

        if (error) {
            enqueueSnackbar(error, {variant: "error"});
        }
    }, [messages, error]);

    if (connectionState === "connecting" || !selectedUser) {
        return (
            <Grid>
                <CircularProgress/>
            </Grid>
        );
    }

    const isMessageFromCurrentUser = (message: WSMessage) => {
        return message.sender_id === parseInt(userID as string);
    };

    const displayMessages = messages.filter(msg =>
        (msg.type === "new_message" || msg.type === "message_history") &&
        msg.content &&
        msg.content !== "Message sent successfully"
    );

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInputText(e.target.value);
    }

    const handleSendMessage = () => {
        if (inputText.trim() === "") {
            return;
        }
        onSendMessage(inputText);
        setInputText("");
    }

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    }

    return (
        <Grid container direction="column" height={"100%"} sx={{p: 1}}>
            <Grid sx={{justifyContent: 'flex-start', alignItems: 'center'}}>
                <Paper sx={{p: 2}}>
                    {selectedUser && (
                        <Grid sx={{display: 'flex', alignItems: 'center', gap: 1}}>
                            <Avatar sx={{bgcolor: getAvatarColor(selectedUser.username), width: 32, height: 32}}>
                                {selectedUser.username.charAt(0).toUpperCase()}
                            </Avatar>
                            <Typography variant="h6" component="h2">
                                {selectedUser.username}
                            </Typography>
                        </Grid>
                    )}
                </Paper>
            </Grid>
            <Grid sx={{flexGrow: 1, overflow: 'auto', mt: 1}}>
                {/*<Paper sx={{p: 2, overflow: 'auto'}}>*/}
                {displayMessages.map((message, index) => {
                    const isFromCurrentUser = isMessageFromCurrentUser(message);
                    return (
                        <Box
                            key={index}
                            sx={{
                                display: 'flex',
                                justifyContent: isFromCurrentUser ? 'flex-end' : 'flex-start',
                                mb: 1,
                            }}
                        >
                            <Box sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: isFromCurrentUser ? 'flex-end' : 'flex-start'
                            }}>
                                <Paper
                                    elevation={1}
                                    sx={{
                                        px: 2,
                                        py: 1,
                                        maxWidth: '70%',
                                        minWidth: 'fit-content',
                                        width: 'auto',
                                        backgroundColor: isFromCurrentUser ? 'primary.main' : 'grey.700',
                                        color: isFromCurrentUser ? 'primary.contrastText' : 'text.primary',
                                        borderRadius: 3,
                                        borderBottomRightRadius: isFromCurrentUser ? 0.5 : 3,
                                        borderBottomLeftRadius: isFromCurrentUser ? 3 : 0.5,
                                        boxShadow: '0 1px 2px rgba(0,0,0,0.1)',
                                        display: 'inline-block',
                                    }}
                                >
                                    <Typography
                                        variant="body2"
                                        sx={{
                                            wordBreak: 'break-word',
                                            overflowWrap: 'break-word',
                                            lineHeight: 1.4,
                                            // whiteSpace: 'pre-wrap'
                                        }}
                                    >
                                        {message.content}
                                    </Typography>
                                </Paper>
                                {message.created_at && (
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            mt: 0.5,
                                            color: 'text.secondary',
                                            fontSize: '0.75rem',
                                        }}
                                    >
                                        {new Date(message.created_at).toLocaleString('en-US', {
                                            hour: 'numeric',
                                            minute: '2-digit',
                                            hour12: true,
                                        })}
                                    </Typography>
                                )}
                            </Box>
                        </Box>
                    );
                })}
                {/*</Paper>*/}
                <div ref={messagesEndRef}/>
            </Grid>
            <Grid>
                <Paper sx={{p: 2}}>
                    <Grid sx={{display: 'flex', alignItems: 'center', gap: 1}}>
                        <TextField
                            fullWidth
                            value={inputText}
                            onChange={handleInputChange}
                            onKeyDown={handleKeyPress}
                            placeholder={`Send a message to ${selectedUser.username}`}
                            disabled={connectionState !== "connected"}
                        />
                        <Tooltip title="Send">
                            <IconButton
                                onClick={handleSendMessage}
                                disabled={connectionState !== "connected"}
                            >
                                <SendIcon/>
                            </IconButton>
                        </Tooltip>
                    </Grid>
                </Paper>
            </Grid>
        </Grid>
    )
}

export default MessageWindow;
