import type {RouteObject} from "react-router-dom";
import {Navigate} from "react-router-dom";
import Register from "../components/Auth/Register.tsx";
import Login from "../components/Auth/Login.tsx";
import ChatPage from "../components/Chat";

const routes: RouteObject[] = [
    {
        path: "/",
        element: <Navigate to="/login" replace/>,
    },
    {
        path: "/login",
        element: <Login/>,
    },
    {
        path: "/register",
        element: <Register/>,
    },
    {
        path: "/chat",
        element: <ChatPage/>,
    }

];

export default routes;
