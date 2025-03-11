import useAuthStore from "../store/authStore.ts";
import { useNavigate } from "react-router";
import Loading from "../components/Loading";
import { useEffect } from "react";
import { GoogleLogo } from "@phosphor-icons/react";

export default function Login() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const login = useAuthStore((state) => state.login);

  useEffect(() => {
    if (user && !loading) {
      navigate(`/member/${user?.username}`);
    }
  }, [user, navigate, loading]);

  if (loading) return <Loading />;

  return (
    <div className="">
      <h1 className="text-2xl font-bold mb-2">Login</h1>
      <button
        onClick={login}
        className="bg-black text-white text-xl rounded-md"
      >
        Login with Google
        <GoogleLogo />
      </button>
    </div>
  );
}
