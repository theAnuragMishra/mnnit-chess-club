import useWebSocketStore from "../store/socketStore"; // Import the Zustand store
import { useParams } from "react-router";
import useChessStore from "../store/gameStore.ts";
import { useState, useEffect } from "react";
import ResultModal from "../components/ResultModal.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { getBaseURL } from "../utils/urlUtils.ts";
import ChessBoard from "../components/ChessBoard.tsx";
import useAuthStore from "../store/authStore.ts";

export default function Play() {
  const params = useParams();
  const { connect } = useWebSocketStore();
  const { setResult, updateFen, updateGround, result, setUserNames } =
    useChessStore();
  const username = useAuthStore((state) => state.user?.username);

  const { data } = useSuspenseQuery({
    queryKey: [params.gameID],
    queryFn: async () => {
      const response = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to fetch game data");
      }
      const x = await response.json();

      setUserNames(x.WhiteUsername, x.BlackUsername);
      updateFen(x.Fen);
      if (x.Result !== "ongoing") {
        setResult(x.Result);
      }
      updateGround();

      return x; // Convert to JSON
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
  const [modalOpen, setModalOpen] = useState(true);

  if (!params.gameID) return <div>Bad Request</div>;
  return (
    <div className="text-2xl px-10 pb-10">
      {result && modalOpen && (
        <ResultModal onClose={() => setModalOpen(false)} />
      )}
      <div className="w-full flex items-center justify-center">
        <div className="mt-5 flex flex-col items-center justify-center">
          <p className="w-full mb-1">
            {data.WhiteUsername !== username
              ? data.WhiteUsername
              : data.BlackUsername}
          </p>
          <ChessBoard gameID={Number(params.gameID)} />
          <p className="w-full mb-1">{username}</p>
        </div>
      </div>
    </div>
  );
}
