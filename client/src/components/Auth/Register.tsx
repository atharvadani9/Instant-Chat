import {useState} from "react";
import {postAPI} from "../../utils/api.ts";
import type {UserAuthRequest} from "../../types";
import {enqueueSnackbar} from "notistack";
import {Alert, Box, Button, Card, CircularProgress, TextField, Typography} from "@mui/material";
import * as React from "react";

export default function Register() {
    const [formData, setFormData] = useState<UserAuthRequest>({
        username: "",
        password: "",
    });
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (formData.username === "" || formData.password === "") {
            setError("Username and password are required");
            return;
        }

        setLoading(true);
        setError("");

        const resp = await postAPI("/user.register", formData);
        if (resp.error && resp.error !== "") {
            enqueueSnackbar(resp.error, {variant: "error"});
        } else {
            localStorage.setItem("userID", resp.user.id.toString());
            enqueueSnackbar("User created successfully", {variant: "success"});
            // @todo: redirect to login page
        }
        setLoading(false);
    };

    return (
        <Box display="flex" justifyContent="center" height="80vh">
            <Card sx={{width: '100%', p: 1}}>
                <Typography variant="h4" component="h1" gutterBottom align="center">
                    {"Register"}
                </Typography>
                <Typography variant="body2" color="text.secondary" align="center" mb={3}>
                    {"Create your account to start chatting"}
                </Typography>
                {error && (
                    <Alert severity="error" sx={{mb: 2}}>
                        {error}
                    </Alert>
                )}
                <Box component={"form"} onSubmit={handleSubmit}
                     sx={{
                         display: 'flex',
                         flexDirection: 'column',
                         alignItems: 'center',
                         gap: 2,
                         width: '100vh'
                     }}>
                    <TextField
                        label="Username"
                        name="username"
                        value={formData.username}
                        onChange={handleChange}
                        margin="normal"
                        required
                        sx={{width: '50%'}}
                    />
                    <TextField
                        label="Password"
                        name="password"
                        type="password"
                        value={formData.password}
                        onChange={handleChange}
                        margin="normal"
                        required
                        sx={{width: '50%'}}
                    />
                    <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        disabled={loading}
                    >
                        {loading ? <CircularProgress/> : "Register"}
                    </Button>
                </Box>
            </Card>
        </Box>
    )
}
