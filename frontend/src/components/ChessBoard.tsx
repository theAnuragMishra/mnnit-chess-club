import React, { useEffect, useRef } from "react";
import { Chessground } from "chessground";
import "../../node_modules/chessground/assets/chessground.base.css";
import "../../node_modules/chessground/assets/chessground.brown.css";
import "../../node_modules/chessground/assets/chessground.cburnett.css";
import useChessStore from "../store/gameStore";
import { getValidMoves } from "../utils/chessUtils";

const ChessBoard: React.FC = () => {
  const boardRef = useRef<HTMLDivElement>(null);
  const chess = useChessStore((state) => state.board);
  const updateHistory = useChessStore((state) => state.updateHistory);
  const setGround = useChessStore((state) => state.setGround);

  useEffect(() => {
    if (!boardRef.current) return;

    // Initialize Chessground
    const ground = Chessground(boardRef.current, {
      draggable: { enabled: true },
      movable: {
        free: false,
        color: "both",
        dests: getValidMoves(chess),
        showDests: true,
        events: {
          after: (orig, dest) => {
            const move = chess.move({ from: orig, to: dest });
            updateHistory(move.san);
            // sendMessage({ type: "move", payload: { MoveStr: moveStr } });
            ground.set({
              fen: chess.fen(),
              turnColor: chess.turn() === "w" ? "white" : "black",
              movable: { dests: getValidMoves(chess) },
            });
          },
        },
      },
      highlight: { lastMove: true, check: true },
    });

    setGround(ground);

    return () => ground.destroy();
  }, [chess]);

  return <div ref={boardRef} className="w-[400px] h-[400px]" />;
};

export default ChessBoard;
