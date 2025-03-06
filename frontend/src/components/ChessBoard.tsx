import { useEffect, useRef } from "react";
import { Chessground } from "chessground";
import "../../node_modules/chessground/assets/chessground.base.css";
import "../../node_modules/chessground/assets/chessground.brown.css";
import "../../node_modules/chessground/assets/chessground.cburnett.css";
import useChessStore from "../store/gameStore";
import { getValidMoves, isPromoting } from "../utils/chessUtils";
import useWebSocketStore from "../store/socketStore";
import useAuthStore from "../store/authStore";
import { Square } from "chess.js";

export default function ChessBoard(props: { gameID: number }) {
  const boardRef = useRef<HTMLDivElement>(null);
  const chess = useChessStore((state) => state.board);
  const check = useChessStore(state=>state.board.isCheck())
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
      // lastMove: [],
      check: check,
      movable: {
        free: false,
        color: white == username ? "white" : "black",
        dests: getValidMoves(chess),
        showDests: true,
        events: {
          after: (orig, dest) => {
            const piece = chess.get(orig as Square);
            if (isPromoting(dest, piece!)) {
              const move = chess.move({ from: orig, to: dest, promotion: "q" });
              sendMessage({
                type: "move",
                payload: {
                  MoveStr: move.san,
                  orig: orig,
                  dest: dest,
                  GameID: Number(props.gameID),
                },
              });
            } else {
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
            }
          },
        },
      },
      highlight: { lastMove: true, check: true },
    });

    setGround(ground);

    return () => ground.destroy();
  }, [chess, props.gameID, sendMessage, setGround, white, username, result,  check]);

  return <div ref={boardRef} className="w-[644px] h-[644px]" />;
}
