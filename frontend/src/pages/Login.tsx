import useAuthStore from "../store/authStore.ts";
import { useNavigate } from "react-router";
import Loading from "../components/Loading";
import { useEffect } from "react";
import { GoogleLogo } from "@phosphor-icons/react";
import {FcGoogle} from "react-icons/fc";

export default function Login() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const login = useAuthStore((state) => state.login);

  useEffect(() => {
    if (!loading && user) {
      navigate(`/member/${user?.username}`);
    }
  }, [user, navigate, loading]);

  if (loading) return <Loading />;

  return (
    <div className="flex flex-col items-center justify-center h-full mt-30 gap-2">
      <h1 className="text-2xl font-bold mb-10">Login</h1>
      <button
        onClick={login}
        className="bg-black px-4 py-2 cursor-pointer text-white text-xl rounded-md flex items-center gap-2"
      >
        Login with Google
        <FcGoogle size={32}/>
      </button>
        <p>More login providers may (or may not) be added.</p>
    </div>
  );
}
