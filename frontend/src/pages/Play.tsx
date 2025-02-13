import { useState } from "react";
import useWebSocketStore from "../store/socketStore"; // Import the Zustand store
import { useParams } from "react-router";
import useChessStore from "../store/gameStore.ts";
import { useEffect } from "react";
import ResultModal from "../components/ResultModal.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { getBaseURL } from "../utils/urlUtils.ts";

export default function Play() {
  const params = useParams();
  const { connect } = useWebSocketStore();
  const [move, setMove] = useState("");
  const { sendMessage } = useWebSocketStore();
  const { player1username, player2username, result } = useChessStore();

  const { data } = useSuspenseQuery({
    queryKey: [params.gameID],
    queryFn: async () => {
      const data = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
        credentials: "include",
      });
      console.log(data);
      return data;
    },
    refetchOnMount: true,
  });

  //     useEffect(() => {
  //     connect(); // Connect WebSocket on mount
  //
  //     return () => close(); // Disconnect on unmount
  // }, [connect, close, sendMessage]);

  useEffect(() => {
    connect(); // Ensure WebSocket stays connected
  }, [connect]);

  if (!params.gameID) return <div>Bad Request</div>;

  const handleSendMessage = () => {
    sendMessage({
      type: "move",
      payload: { MoveStr: move, GameID: Number(params.gameID!) },
    });
    setMove("");
  };

  return (
    <div>
      {result && <ResultModal />}
      <h2>Play</h2>
      <div>
        White: {player1username} Black: {player2username}
      </div>
      <div>Game ID: {params.gameID}</div>
      Move:
      <input
        type="text"
        className="w-full px-3 py-2 border rounded"
        value={move}
        onChange={(e) => setMove(e.target.value)}
      />
      <button
        onClick={handleSendMessage}
        className="rounded-md border mt-4 p-4  "
      >
        Send
      </button>
    </div>
  );
}
