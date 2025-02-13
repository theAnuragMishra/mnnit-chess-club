import useAuthStore from "../store/authStore.ts";
import Loading from "../components/Loading";
import useWebSocketStore from "../store/socketStore.ts";
import useChessStore from "../store/gameStore.ts";
import { useEffect } from "react";
import { useNavigate } from "react-router";
import { getBaseURL } from "../utils/urlUtils.ts";

export default function Home() {
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);

  const { connect } = useWebSocketStore();
  const { gameID } = useChessStore();
  const navigate = useNavigate();

  useEffect(() => {
    if (gameID) {
      navigate(`/play/${gameID}`);
    }
  }, [gameID, navigate]);

  useEffect(() => {
    connect();
  }, [connect]);

  async function handleInitGame() {
    await fetch(`${getBaseURL()}/game/init`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ username: user?.username }),
    });
  }

  if (loading) return <Loading />;

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
