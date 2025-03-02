import { useEffect, useRef } from "react";
import { Chessground } from "chessground";
import "../../node_modules/chessground/assets/chessground.base.css";
import "../../node_modules/chessground/assets/chessground.brown.css";
import "../../node_modules/chessground/assets/chessground.cburnett.css";
import useChessStore from "../store/gameStore";
import { getValidMoves } from "../utils/chessUtils";
import useWebSocketStore from "../store/socketStore";
import useAuthStore from "../store/authStore";

export default function ChessBoard(props: { gameID: number }) {
  const boardRef = useRef<HTMLDivElement>(null);
  const chess = useChessStore((state) => state.board);
  const setGround = useChessStore((state) => state.setGround);
  const sendMessage = useWebSocketStore((state) => state.sendMessage);
  const white = useChessStore((state) => state.whiteusername);
  const username = useAuthStore((state) => state.user?.username);
  const result = useChessStore((state) => state.result);

  useEffect(() => {
    if (!boardRef.current) return;

    // Initialize Chessground
    const ground = Chessground(boardRef.current, {
      fen: chess.fen(),
      orientation: white == username ? "white" : "black",
      draggable: { enabled: true },
      turnColor: chess.turn() == "w" ? "white" : "black",
      viewOnly: result !== "" && result !== "ongoing",
      lastMove: [],
      movable: {
        free: false,
        color: white == username ? "white" : "black",
        dests: getValidMoves(chess),
        showDests: true,
        events: {
          after: (orig, dest) => {
            const move = chess.move({ from: orig, to: dest });
            sendMessage({
              type: "move",
              payload: {
                MoveStr: move.san,
                orig: orig,
                dest: dest,
                GameID: Number(props.gameID),
              },
            });
          },
        },
      },
      highlight: { lastMove: true, check: true },
    });

    setGround(ground);

    return () => ground.destroy();
  }, [chess, props.gameID, sendMessage, setGround, white, username, result]);

  return <div ref={boardRef} className="w-[644px] h-[644px]" />;
}
