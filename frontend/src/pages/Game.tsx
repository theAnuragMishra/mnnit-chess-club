import useWebSocketStore from "../store/socketStore.ts"; // Import the Zustand store
import { useParams } from "react-router";
import useChessStore from "../store/gameStore.ts";
import { useState, useEffect } from "react";
import ResultModal from "../components/ResultModal.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { getBaseURL } from "../utils/urlUtils.ts";
import ChessBoard from "../components/ChessBoard.tsx";
import useAuthStore from "../store/authStore.ts";
import { chunkArray } from "../utils/utils.ts";
import Clock from "../components/Clock.tsx";

export default function Game() {
  console.log("on game page");
  const params = useParams();
  const { connect } = useWebSocketStore();
  const {
    setResult,
    updateFen,
    updateGround,
    result,
    board,
    timeBlack,
    timeWhite,
    moveHistory,
    setUserNames,
    setHistory,
    setTimeBlack,
    setTimeWhite,
    setIncrement,
  } = useChessStore();
  const username = useAuthStore((state) => state.user?.username);

  const { data } = useSuspenseQuery({
    queryKey: [params.gameID, "gameinfo"],
    queryFn: async () => {
      const response = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Failed to fetch game data");
      }
      const x = await response.json();
      console.log(x);

      setUserNames(x.game.WhiteUsername, x.game.BlackUsername);
      setHistory(x.moves);
      updateFen(x.game.Fen);
      setTimeWhite(x.timeWhite);
      setTimeBlack(x.timeBlack);
      setIncrement(x.game.Increment);
      // setHistory(x.moves);
      if (x.game.Result !== "ongoing") {
        setResult(x.game.Result);
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

  const history = moveHistory ? chunkArray(moveHistory, 2) : [];

  return (
    <div className="text-2xl px-10 pb-10">
      {result && modalOpen && (
        <ResultModal onClose={() => setModalOpen(false)} />
      )}
      <div className="w-full flex gap-15 items-center justify-center">
        <div className="mt-5 flex flex-col items-center justify-center">
          <p className="w-full mb-1">
            {data.game.WhiteUsername === username
              ? data.game.BlackUsername
              : data.game.WhiteUsername}
          </p>
          <ChessBoard gameID={Number(params.gameID)} />
          <p className="w-full mb-1">{username}</p>
        </div>
        <div className="w-[200px] h-full">
          {history &&
            history.map((move, index) => {
              return (
                <p key={index} className="flex w-full justify-between">
                  <span>
                    {index + 1}
                    {".    "}
                    {move[0] && move[0].MoveNotation}
                  </span>
                  <span> {move[1] && move[1].MoveNotation}</span>
                </p>
              );
            })}
        </div>
        <div>
          <Clock
            initialWhite={timeWhite * 1000}
            initialBlack={timeBlack * 1000}
            onTimeUp={() => alert("timeup")}
            turn={board.turn() === "w" ? "white" : "black"}
          />
        </div>
      </div>
    </div>
  );
}
