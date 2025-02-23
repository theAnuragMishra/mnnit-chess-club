import useAuthStore from "../store/authStore.ts";
import Loading from "../components/Loading";
import { useNavigate } from "react-router";

export default function Home() {
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const navigate = useNavigate();
  if (loading) return <Loading />;

  const handleInitGame = () => {
    navigate("/play");
  };

  return (
    <div className="flex flex-col items-center justify-center h-full ">
      {user ? (
        <button
          onClick={handleInitGame}
          className="text-4xl bg-gray-500 p-3 rounded-lg cursor-pointer"
        >
          Play!
        </button>
      ) : (
        <p>Login to play</p>
      )}
    </div>
  );
}
