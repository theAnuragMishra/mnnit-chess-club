import Nav from "../components/nav.tsx";
import {Outlet, useNavigate} from "react-router";
import {useEffect} from "react";
import useAuthStore from "../store/authStore.ts";

export default function RootLayout() {
    const navigate = useNavigate();
    const user = useAuthStore((state) => state.user);
    const loading = useAuthStore((state) => state.loading);
    useEffect(() => {
        if (!loading && user && !user?.username) {
            navigate("/set-username");
        }
    });
    return (
        <main className="min-h-screen flex-col bg-gray-900 text-white items-center">
            <Nav />
            <Outlet />
        </main>
    );
}
