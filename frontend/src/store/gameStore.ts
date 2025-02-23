import { create } from "zustand";
import { Chess } from "chess.js";
import { Api } from "chessground/api";
import { getValidMoves } from "../utils/chessUtils";

interface ChessState {
  board: Chess;
  whiteusername: string;
  blackusername: string;
  timeBlack: number;
  timeWhite: number;
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
  setTimeBlack: (time: number) => void;
  setTimeWhite: (time: number) => void;
}

const useChessStore = create<ChessState>()((set, get) => ({
  ground: null,
  whiteusername: "",
  blackusername: "",
  timeBlack: 0,
  timeWhite: 0,
  result: "",
  board: new Chess(),
  moveHistory: [],

  setUserNames: (white: string, black: string) =>
    set({ whiteusername: white, blackusername: black }),

  setGround: (ground: Api) => set({ ground }),

  updateGround: () => {
    get().ground?.set({
      // viewOnly: get().result !== "" && get().result !== "ongoing",
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
  setTimeWhite: (time: number) => set({ timeWhite: time }),
  setTimeBlack: (time: number) => set({ timeBlack: time }),

  resetGame: () =>
    set(() => ({
      board: new Chess(),
      moveHistory: [],
    })),
}));

export default useChessStore;
