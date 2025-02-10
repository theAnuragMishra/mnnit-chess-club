import { create } from "zustand";
import {getBaseURL} from "../utils/urlUtils.ts";

interface User{
    username:string;
    userID: string;
}

interface AuthState{
    user: User|null;
    login: (username: string, password: string) => Promise<void>;
    register: (username: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    checkAuth: ()=> Promise <void>
    loading: boolean
}

// Create Zustand store
const useAuthStore = create<AuthState>()((set) => ({
    user: null,
    loading: true,
    login: async (username:string, password:string) => {

        const res = await fetch(`${getBaseURL()}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ username, password }),
        });

        if (!res.ok) {throw new Error("Login failed");}

        const user = await res.json();
        // console.log(user);
        set({ user, loading:false });
    },
    register: async (username: string, password: string) => {

            const res = await fetch(`${getBaseURL()}/register`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
            });

            if (!res.ok) throw new Error("Registration failed");

    },
    logout: async () => {
        await fetch(`${getBaseURL()}/logout`, {
            method: "POST",
            credentials: "include",
        });

        set({ user: null, loading: false });
    },
    checkAuth: async () => {
set({loading:true});
        // await delay(10000);

            const res = await fetch(`${getBaseURL()}/me`, { credentials: "include" });


if(!res.ok){
    set({ user: null, loading: false });
    return;
}
        const user = await res.json();
        set({ user, loading: false });

    },
}));



export default useAuthStore;
