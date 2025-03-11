import { create } from "zustand";
import { getBaseURL } from "../utils/urlUtils.ts";
import useWebSocketStore from "./socketStore.ts";

interface User {
  username: string | null;
  userID: string;
}

interface AuthState {
  user: User | null;
  login: () => void;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  loading: boolean;
}

// Create Zustand store
const useAuthStore = create<AuthState>()((set, get) => ({
  user: null,
  loading: true,
  login: () => {
    window.location.href = "http://localhost:8080/auth/login/google";
  },

  logout: async () => {
    await fetch(`${getBaseURL()}/logout`, {
      method: "POST",
      credentials: "include",
    });

    set({ user: null, loading: false });
    useWebSocketStore.getState().close();
  },
  checkAuth: async () => {
    set({ loading: true });
    // await delay(10000);

    const res = await fetch(`${getBaseURL()}/me`, { credentials: "include" });

    if (!res.ok) {
      set({ user: null, loading: false });
      return;
    }
    const user = await res.json();
    set({ user, loading: false });
  },
}));

export default useAuthStore;
