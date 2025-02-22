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

  async function handleInitGame(timerCode: number) {
    await fetch(`${getBaseURL()}/game/init`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ username: user?.username, timerCode }),
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
        onClick={() => handleInitGame(2)}
        className="w-[150px] h-[100px] bg-gray-400 cursor-pointer text-xl"
      >
        3+2
      </button>
      <button
        onClick={() => handleInitGame(3)}
        className="w-[150px] h-[100px] bg-gray-400 cursor-pointer text-xl"
      >
        10+0
      </button>
      <button
        onClick={() => handleInitGame(1)}
        className="w-[150px] h-[100px] bg-gray-400 cursor-pointer text-xl"
      >
        1+0
      </button>
    </div>
  );
}
