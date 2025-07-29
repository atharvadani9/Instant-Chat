import type {RouteObject} from "react-router-dom";
import Register from "../components/Auth/Register.tsx";

const routes: RouteObject[] = [
    {
        path: "/register",
        element: <Register/>,
    },
    {
        path: "*",
        element: <></>,
    },
];

export default routes;
