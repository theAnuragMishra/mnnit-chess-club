import Nav from "../components/nav.tsx";
import {Outlet} from "react-router";

export default function RootLayout() {

    return (
    <main className="min-h-screen flex-col bg-gray-700 text-white">
        <Nav/>
        <Outlet/>
    </main>
)
}