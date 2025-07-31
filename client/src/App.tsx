import {SnackbarProvider, closeSnackbar} from "notistack";
import './App.css'
import {useRoutes} from "react-router-dom";
import routes from "./routes/router";
import {CssBaseline, ThemeProvider} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import IconButton from "@mui/material/IconButton";
import darkTheme from "./utils/theme.ts";
import {AuthProvider} from "./contexts/AuthContext.tsx";


function App() {
    const content = useRoutes(routes)

    return (
        <AuthProvider>
            <ThemeProvider theme={darkTheme}>
                <CssBaseline/>
                <SnackbarProvider
                    maxSnack={3}
                    anchorOrigin={{
                        vertical: "top",
                        horizontal: "right",
                    }}
                    autoHideDuration={4000}
                    action={(snackbarId) => (
                        <IconButton
                            size="small"
                            aria-label="close"
                            color="inherit"
                            onClick={() => closeSnackbar(snackbarId)}
                        >
                            <CloseIcon fontSize="small"/>
                        </IconButton>
                    )}>
                    {content}
                </SnackbarProvider>
            </ThemeProvider>
        </AuthProvider>
    )
}

export default App
