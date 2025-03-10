import { Key } from "chessground/types";
import { Chess, Color, Piece } from "chess.js";

export function colorToCgColor(chessjsColor: Color): "white" | "black" {
  return chessjsColor === "w" ? "white" : "black";
}
export function cgColorToColor(chessgroundColor: "white" | "black"): Color {
  return chessgroundColor === "white" ? "w" : "b";
}

export const getValidMoves = (chess: Chess) => {
  const moves = new Map<Key, Key[]>();

  chess.board().forEach((row) => {
    row.forEach((square) => {
      if (square) {
        const from = square.square as Key;
        const legalMoves = chess
          .moves({ square: square.square, verbose: true })
          .map((m) => m.to);
        if (legalMoves.length) moves.set(from, legalMoves);
      }
    });
  });

  return moves;
};

export function isPromoting(dest: Key, piece: Piece) {
  return (
    piece.type == "p" &&
    ((piece.color == "w" && dest[1] == "8") ||
      (piece.color == "b" && dest[1] == "1"))
  );
}
