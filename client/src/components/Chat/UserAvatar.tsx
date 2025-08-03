// noinspection JSIgnoredPromiseFromCall

import React, {useEffect, useState} from 'react';
import {Avatar, Box, IconButton, Menu, MenuItem, Typography} from '@mui/material';
import {useNavigate} from 'react-router-dom';
import {useAuthContext} from '../../contexts/AuthContext';
import {getAPI} from '../../utils/api';
import {enqueueSnackbar} from 'notistack';
import type {User} from '../../types';
import {getAvatarColor} from '../../utils/utils';

const UserAvatar: React.FC = () => {
    const {userID, logout} = useAuthContext();
    const navigate = useNavigate();
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    useEffect(() => {
        const fetchCurrentUser = async () => {
            if (!userID) return;

            try {
                const response = await getAPI(`/user.get.me?user_id=${userID}`);
                if (response.error) {
                    enqueueSnackbar(response.error, {variant: "error"});
                } else {
                    setCurrentUser(response.user);
                }
            } catch (err) {
                enqueueSnackbar('Failed to fetch user data', {variant: "error"});
            }
        };

        fetchCurrentUser();
    }, [userID]);

    const handleAvatarClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
    };

    const handleSignOut = () => {
        logout();
        navigate("/login");
        handleMenuClose();
    };

    if (!currentUser) {
        return null;
    }

    return (
        <>
            <IconButton onClick={handleAvatarClick} size="small">
                <Avatar
                    sx={{
                        bgcolor: getAvatarColor(currentUser.username),
                        width: 40,
                        height: 40
                    }}
                >
                    {currentUser.username.charAt(0).toUpperCase()}
                </Avatar>
            </IconButton>
            <Menu
                anchorEl={anchorEl}
                open={open}
                onClose={handleMenuClose}
                onClick={handleMenuClose}
                transformOrigin={{horizontal: 'right', vertical: 'top'}}
                anchorOrigin={{horizontal: 'right', vertical: 'bottom'}}
                sx={{
                    '& .MuiPaper-root': {
                        minWidth: 200,
                        mt: 1
                    }
                }}
            >
                <Box sx={{px: 2, py: 1, borderBottom: 1, borderColor: 'divider'}}>
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1}}>
                        <Avatar
                            sx={{
                                bgcolor: getAvatarColor(currentUser.username),
                                width: 32,
                                height: 32
                            }}
                        >
                            {currentUser.username.charAt(0).toUpperCase()}
                        </Avatar>
                        <Typography variant="body1" fontWeight="medium">
                            {currentUser.username}
                        </Typography>
                    </Box>
                </Box>
                <MenuItem onClick={handleSignOut}>
                    {"Sign out"}
                </MenuItem>
            </Menu>
        </>
    );
};

export default UserAvatar;
