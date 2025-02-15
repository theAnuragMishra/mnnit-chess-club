import { create } from "zustand";
import { Chess } from "chess.js";
import { Api } from "chessground/api";
import { getValidMoves } from "../utils/chessUtils";

interface ChessState {
  board: Chess;
  ground: Api | null;
  result: string | null;
  gameID: string;
  player1username: string;
  player2username: string;
  moveHistory: string[];
  updateHistory: (move: string) => void;
  updateGround: () => void;
  setGround: (ground: Api) => void;
  resetGame: () => void;
}

const useChessStore = create<ChessState>((set, get) => ({
  player1username: "",
  player2username: "",
  ground: null,
  gameID: "",
  result: null,
  board: new Chess(),
  moveHistory: [],

  setGround: (ground: Api) => set({ ground }),

  updateGround: () => {
    get().ground?.set({
      fen: get().board.fen(),
      turnColor: get().board.turn() === "w" ? "white" : "black",
      movable: { dests: getValidMoves(get().board) },
    });
  },

  updateHistory: (move: string) => {
    set((state) => ({ moveHistory: [...state.moveHistory, move] }));
    console.log(get().moveHistory);
  },
  resetGame: () =>
    set(() => ({
      board: new Chess(),
      moveHistory: [],
    })),
}));

export default useChessStore;
