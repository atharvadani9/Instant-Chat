import {useState} from "react";
import {postAPI} from "../../utils/api.ts";
import type {UserAuthRequest} from "../../types";
import {enqueueSnackbar} from "notistack";
import {Alert, Box, Button, CircularProgress, TextField, Typography, Grid} from "@mui/material";
import * as React from "react";
import {Link} from "react-router-dom";

export default function Login() {
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

        const resp = await postAPI("/user.login", formData);
        if (resp.error && resp.error !== "") {
            enqueueSnackbar(resp.error, {variant: "error"});
        } else {
            localStorage.setItem("userID", resp.user.id.toString());
            enqueueSnackbar("Login successful", {variant: "success"});
            // @todo: redirect to chat page
        }
        setLoading(false);
    };

    return (
        <Grid
            display="flex"
            flexDirection="column"
            justifyContent="flex-start"
            alignItems="center"
            sx={{minHeight: '100vh', width: '100%', pt: 4, gap: 1}}
        >
            <Typography variant="h4" component="h1" gutterBottom align="center">
                {"Login"}
            </Typography>
            <Typography variant="body2" color="text.secondary" align="center">
                {"Welcome back! Please sign in to your account"}
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
                     gap: 1,
                     width: '100%',
                     maxWidth: 400
                 }}>
                <TextField
                    label="Username"
                    name="username"
                    value={formData.username}
                    onChange={handleChange}
                    margin="normal"
                    required
                    sx={{width: '100%'}}
                />
                <TextField
                    label="Password"
                    name="password"
                    type="password"
                    value={formData.password}
                    onChange={handleChange}
                    margin="normal"
                    required
                    sx={{width: '100%'}}
                />
                <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    disabled={loading}
                >
                    {loading ? <CircularProgress/> : "Login"}
                </Button>
                <Link to="/register" style={{paddingTop: 2}}>
                    {"New user? Go to Register"}
                </Link>
            </Box>
        </Grid>
    )
}