import React from "react";
import type {User} from "../../types";
import {enqueueSnackbar} from "notistack";
import {Avatar, CircularProgress, Grid, ListItemAvatar, ListItemButton, ListItemText} from "@mui/material";
import {getAvatarColor} from "../../utils/utils.ts";

type UserListProps = {
    users: User[];
    selectedUser: User | null;
    onUserSelect: (user: User) => void;
    loading: boolean;
    error: string | null;
}

const UserList: React.FC<UserListProps> = ({users, selectedUser, onUserSelect, loading, error}) => {
    if (error) {
        enqueueSnackbar(error, {variant: "error"});
        return <></>;
    }

    return (
        <>
            {loading ? (
                <Grid>
                    <CircularProgress/>
                </Grid>
            ) : (
                <Grid container direction="column">
                    {users.map((user) => (
                        <Grid key={user.id}>
                            <ListItemButton
                                selected={selectedUser?.id === user.id}
                                onClick={() => onUserSelect(user)}
                            >
                                <ListItemAvatar>
                                    <Avatar sx={{bgcolor: getAvatarColor(user.username)}}>
                                        {user.username.charAt(0).toUpperCase()}
                                    </Avatar>
                                </ListItemAvatar>
                                <ListItemText primary={user.username}/>
                            </ListItemButton>
                        </Grid>
                    ))}
                </Grid>
            )}
        </>
    )
}

export default UserList;
