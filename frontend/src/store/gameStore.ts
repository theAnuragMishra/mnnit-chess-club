import { create } from "zustand";
import { Chess } from "chess.js";
import { Api } from "chessground/api";
import { getValidMoves } from "../utils/chessUtils";

interface ChessState {
  board: Chess;
  whiteusername: string;
  blackusername: string;
  ground: Api | null;
  result: string;
  moveHistory: { MoveFen: string; MoveNotation: string; MoveNumber: number }[];
  updateFen: (fen: string) => void;
  setResult: (result: string) => void;
  updateHistory: (move: string) => void;
  setHistory: (his: string[]) => void;
  updateGround: () => void;
  setGround: (ground: Api) => void;
  resetGame: () => void;
  setUserNames: (white: string, black: string) => void;
}

const useChessStore = create<ChessState>()((set, get) => ({
  ground: null,
  whiteusername: "",
  blackusername: "",
  result: "",
  board: new Chess(),
  moveHistory: [],

  setUserNames: (white: string, black: string) =>
    set({ whiteusername: white, blackusername: black }),

  setGround: (ground: Api) => set({ ground }),

  updateGround: () => {
    get().ground?.set({
      // viewOnly: get().result !== "",
      fen: get().board.fen(),
      turnColor: get().board.turn() === "w" ? "white" : "black",
      movable: { dests: getValidMoves(get().board) },
    });
  },

  setResult: (result: string) => set({ result }),

  updateFen: (fen: string) => {
    get().board.load(fen);
  },

  setHistory: (his: any) => {
    get().moveHistory = his;
  },

  updateHistory: (move: any) => {
    if (get().moveHistory)
      set((state) => ({ moveHistory: [...state.moveHistory, move] }));
    else {
      set(() => ({ moveHistory: [move] }));
    }
  },
  resetGame: () =>
    set(() => ({
      board: new Chess(),
      moveHistory: [],
    })),
}));

export default useChessStore;
