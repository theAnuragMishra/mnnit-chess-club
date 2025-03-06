import { create } from "zustand";
import { Chess } from "chess.js";
import { Api } from "chessground/api";
import { getValidMoves } from "../utils/chessUtils";
import { Key } from "chessground/types";

interface ChessState {
  board: Chess;
  whiteusername: string;
  blackusername: string;
  whiteID: string;
  blackID: string;
  timeBlack: number;
  timeWhite: number;
  ground: Api | null;
  result: string;
  moveHistory:
  | { MoveFen: string; MoveNotation: string; MoveNumber: number }[]
  | null;
  activeIndex: number;
  updateFen: (fen: string) => void;
  setResult: (result: string) => void;
  updateHistory: (move: string) => void;
  setHistory: (his: string[]) => void;
  updateGround: () => void;
  makeMoveOnGround: (s1: Key, s2: Key) => void;
  setGroundFen: (fen: string) => void;
  setGroundViewOnly: (x: boolean) => void;
  setGroundLastMoves: (orig: string, dest: string) => void;
  setGround: (ground: Api) => void;
  resetGame: () => void;
  setUserNames: (white: string, black: string) => void;
  setUserIDs: (white: string, black: string) => void;
  setTimeBlack: (time: number) => void;
  setTimeWhite: (time: number) => void;
  setActiveIndex: (index: number) => void;
}

const useChessStore = create<ChessState>()((set, get) => ({
  ground: null,
  whiteusername: "",
  blackusername: "",
  whiteID: "",
  blackID: "",
  timeBlack: 0,
  timeWhite: 0,
  result: "",
  board: new Chess(),
  moveHistory: null,
  activeIndex: -1,

  setActiveIndex: (index: number) => set({ activeIndex: index }),

  setUserNames: (white: string, black: string) =>
    set({ whiteusername: white, blackusername: black }),

  setUserIDs: (white: string, black: string) =>
    set({ whiteID: white, blackID: black }),

  setGround: (ground: Api) => set({ ground }),

  makeMoveOnGround: (s1: Key, s2: Key) => {
    get().ground?.move(s1, s2);
  },

  setGroundFen: (fen: string) => {
    get().ground?.set({
      fen: fen,
    });
  },

  setGroundLastMoves: (orig: string, dest: string) => {
    get().ground?.set({
      lastMove: [orig as Key, dest as Key],
    });
  },

  setGroundViewOnly: (x: boolean) => {
    get().ground?.set({
      viewOnly: x,
    });
  },

  updateGround: () => {
    get().ground?.set({
      turnColor: get().board.turn() === "w" ? "white" : "black",
      movable: { dests: getValidMoves(get().board) },
    });
  },

  setResult: (result: string) => set({ result }),

  updateFen: (fen: string) => {
    get().board.load(fen);
    // const newBoard = get().board ; // Clone or create a new board instance
    // newBoard.load(fen);
    // set({ board: newBoard }); // Now Zustand detects a change
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
