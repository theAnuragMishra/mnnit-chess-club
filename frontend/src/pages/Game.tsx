import useWebSocketStore from "../store/socketStore.ts"; // Import the Zustand store
import { useParams } from "react-router";
import useChessStore from "../store/gameStore.ts";
import { useState, useEffect } from "react";
import ResultModal from "../components/ResultModal.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { getBaseURL } from "../utils/urlUtils.ts";
import ChessBoard from "../components/ChessBoard.tsx";
import useAuthStore from "../store/authStore.ts";
import Clock from "../components/Clock2.tsx";
import HistoryTable from "../components/HistoryTable.tsx";
import Chat from "../components/Chat.tsx";
import GameInfo from "../components/GameInfo.tsx";

export default function Game() {
  const params = useParams();
  const { connect, sendMessage } = useWebSocketStore();
  const {
    setResult,
    setUserIDs,
    setGroundLastMoves,
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
    setActiveIndex,
    setReason,
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

      setUserNames(x.game.WhiteUsername, x.game.BlackUsername);
      setUserIDs(x.game.WhiteID, x.game.BlackID);
      setHistory(x.moves);
      updateFen(x.game.Fen);
      setTimeWhite(x.timeWhite);
      setTimeBlack(x.timeBlack);
      setResult(x.game.Result);
      if (x.game.Result != "ongoing" || x.game.Result != "") {
        setReason(x.game.ResultReason);
      }
      updateGround();
      if (x.moves) {
        setActiveIndex(x.moves.length - 1);
        setGroundLastMoves(
          x.moves[x.moves.length - 1].Orig,
          x.moves[x.moves.length - 1].Dest,
        );
      }
      // console.log(x);
      return x; // Convert to JSON
    },
    refetchOnMount: true,
  });

  useEffect(() => {
    connect(); // Ensure WebSocket stays connected
  }, [connect]);
  const [modalOpen, setModalOpen] = useState(true);

  if (!params.gameID) return <div>Bad Request</div>;

  const whiteUp = data.game.WhiteUsername !== username;

  return (
    <div className="text-2xl px-5">
      {result !== "ongoing" && result !== "" && modalOpen && (
        <ResultModal onClose={() => setModalOpen(false)} />
      )}
      <div className="w-full flex items-center justify-around">
        {result !== "1-0" && result !== "0-1" ? (
          <Chat gameID={params.gameID} />
        ) : (
          <GameInfo />
        )}
        <div className="flex items-center justify-center">
          <ChessBoard gameID={Number(params.gameID)} />
        </div>
        <div className="w-1/4 h-full flex flex-col gap-2">
          <p className="w-full mb-1 flex justify-between items-center">
            {whiteUp ? data.game.WhiteUsername : data.game.BlackUsername}{" "}
            <Clock
              initialTime={whiteUp ? timeWhite : timeBlack}
              onTimeUp={() => {
                // console.log("time up");
                sendMessage({
                  type: "timeup",
                  payload: {
                    color: whiteUp ? "black" : "white",
                    gameID: Number(params.gameID),
                  },
                });
              }}
              active={
                result !== "ongoing" && result !== ""
                  ? false
                  : whiteUp
                    ? board.turn() === "w"
                    : board.turn() === "b"
              }
            />
          </p>

          <HistoryTable history={moveHistory} />

          <p className="w-full mb-1 flex items-center justify-between">
            {username}{" "}
            <Clock
              initialTime={whiteUp ? timeBlack : timeWhite}
              onTimeUp={() => {
                // console.log("time up");
                sendMessage({
                  type: "timeup",
                  payload: {
                    color: whiteUp ? "black" : "white",
                    gameID: Number(params.gameID),
                  },
                });
              }}
              active={
                result !== "ongoing" && result !== ""
                  ? false
                  : whiteUp
                    ? board.turn() === "b"
                    : board.turn() === "w"
              }
            />
          </p>
        </div>
      </div>
    </div>
  );
}
