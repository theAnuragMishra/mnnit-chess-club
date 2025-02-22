import { useNavigate } from "react-router";
import useWebSocketStore from "../store/socketStore";
import { useEffect } from "react";
import { getBaseURL } from "../utils/urlUtils";
import useAuthStore from "../store/authStore";

export default function Play() {
  console.log("on play page");
  const { connect } = useWebSocketStore();
  const { setNavigate } = useWebSocketStore();
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);

  useEffect(() => {
    setNavigate(navigate);
    connect();
  }, [connect, navigate, setNavigate]);

  async function handleInitGame() {
    await fetch(`${getBaseURL()}/game/init`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ username: user?.username }),
    });
  }
  // useEffect(() => {
  //
  //   handleInitGame();
  // }, [user]);
  //
  return (
    <div className="w-full min-h-screen flex items-center justify-center">
      <button
        onClick={handleInitGame}
        className="w-[100px] h-[500px] bg-gray-400 cursor-pointer text-xl"
      >
        Start
      </button>
    </div>
  );
}
