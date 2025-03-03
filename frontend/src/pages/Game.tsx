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
import Clock from "../components/Clock2.tsx";

export default function Game() {
  const params = useParams();
  const { connect, sendMessage } = useWebSocketStore();
  const {
    setResult,
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
    setGroundFen,
    setGroundViewOnly,
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
      setHistory(x.moves);
      updateFen(x.game.Fen);
      setTimeWhite(x.timeWhite);
      setTimeBlack(x.timeBlack);
      setResult(x.game.Result);
      updateGround();
      if (x.moves) {
        setGroundLastMoves(
          x.moves[x.moves.length - 1].Orig,
          x.moves[x.moves.length - 1].Dest,
        );
      }

      return x; // Convert to JSON
    },
    refetchOnMount: true,
  });

  useEffect(() => {
    connect(); // Ensure WebSocket stays connected
  }, [connect]);
  const [modalOpen, setModalOpen] = useState(true);

  if (!params.gameID) return <div>Bad Request</div>;

  const history = moveHistory ? chunkArray(moveHistory, 2) : [];

  const whiteUp = data.game.WhiteUsername !== username;

  return (
    <div className="text-2xl px-5">
      {result !== "ongoing" && result !== "" && modalOpen && (
        <ResultModal onClose={() => setModalOpen(false)} />
      )}
      <div className="w-full flex items-center justify-around">
        <div className="w-1/4 px-2 py-1 bg-white text-black">
          <h1>Chat</h1>
          <ul>
            <li>hi this</li>
            <li>hello will</li>
            <li>there be</li>
            <li>is replaced</li>
            <li>me by</li>
            <li>bye chat</li>
          </ul>
        </div>
        <div className="flex items-center justify-center">
          <ChessBoard gameID={Number(params.gameID)} />
        </div>
        <div className="w-1/4 h-full flex flex-col gap-2">
          <p className="w-full mb-1 flex justify-between items-center">
            {whiteUp ? data.game.WhiteUsername : data.game.BlackUsername}{" "}
            <Clock
              initialTime={whiteUp ? timeWhite : timeBlack}
              onTimeUp={() => {
                console.log("time up");
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
          <div className="h-[300px] text-lg px-4 py-2  bg-gray-800 relative">
            <div className="overflow-y-auto h-[250px] overflow-x-hidden">
              {history &&
                history.map((move, index) => {
                  return (
                    <div key={index} className="w-full grid grid-cols-[1fr_16fr_16fr] gap-10">
                      <span>
                        {index + 1}
                        {".    "}
                      </span>
                      {move[0] && (
                        <button
                          onClick={() => {
                            setGroundFen(move[0].MoveFen);
                            setGroundLastMoves(move[0].Orig, move[0].Dest);
                            if (index === history.length - 1 && !move[1]) {
                              setGroundViewOnly(false);
                            } else {
                              setGroundViewOnly(true);
                            }
                          }}
                          className="cursor-pointer"
                        >
                          {move[0].MoveNotation}
                        </button>
                      )}
                      {move[1] && (
                        <button
                          onClick={() => {
                            setGroundFen(move[1].MoveFen);
                            setGroundLastMoves(move[1].Orig, move[1].Dest);
                            if (index === history.length - 1) {
                              setGroundViewOnly(false);
                            } else {
                              setGroundViewOnly(true);
                            }
                          }}
                          className="cursor-pointer"
                        >
                          {move[1].MoveNotation}
                        </button>
                      )}
                    </div>
                  );
                })}
            </div>
            <div className="absolute bottom-2 w-full flex justify-around">
              <button
                className="cursor-pointer"
                onClick={() => {
                  setGroundFen(history[0][0].MoveFen);
                  setGroundLastMoves(history[0][0].Orig, history[0][0].Dest);
                  setGroundViewOnly(true);
                }}
              >
                {"<"}
              </button>
              <button
                className="cursor-pointer"
                onClick={() => {
                  setGroundFen(board.fen());
                  if (history[history.length - 1][1]) {
                    setGroundLastMoves(
                      history[history.length - 1][1].Orig,
                      history[history.length - 1][1].Dest,
                    );
                  } else {
                    setGroundLastMoves(
                      history[history.length - 1][0].Orig,
                      history[history.length - 1][0].Dest,
                    );
                  }
                  setGroundViewOnly(false);
                }}
              >
                {">"}
              </button>
            </div>
          </div>
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
