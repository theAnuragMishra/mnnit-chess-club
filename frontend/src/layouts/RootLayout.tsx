import Nav from "../components/nav.tsx";
import { Outlet } from "react-router";

export default function RootLayout() {
    return (
        <main className="min-h-screen flex-col bg-gray-900 text-white items-center">
            <Nav />
            <Outlet />
        </main>
    );
}
