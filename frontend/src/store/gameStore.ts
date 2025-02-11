import { create } from "zustand";
import { Chess } from "chess.js";

interface ChessState {
    board: Chess;
    result: string | null;
    gameID: string;
    player1username: string;
    player2username: string;
    turn: "white" | "black";
    moveHistory: string[];
    makeMove: (move: string) => boolean;  // Returns true if move is valid
    resetGame: () => void;
}

const useChessStore = create<ChessState>((set, get) => ({
    player1username: "",
    player2username: "",
    gameID: "",
    result: null,
    board: new Chess(), // Keep a single Chess instance
    turn: "white",
    moveHistory: [],

    setResult: (result:string) => set({result: result}),
    makeMove: (move: string) => {
        const currentBoard = get().board;  // Get the current board (mutable)
        const moveResult = currentBoard.move(move); // Try making the move

        if (moveResult) {
            set((state) => ({
                board: currentBoard, // Reuse existing instance
                turn: state.turn === "white" ? "black" : "white",
                moveHistory: [...state.moveHistory, move],
            }));
            return true;
        } else {
            console.warn("Invalid move:", move);
            return false;
        }
    },

    resetGame: () => set(() => ({
        board: new Chess(), // Reset with a new instance only when needed
        turn: "white",
        moveHistory: [],
    })),
}));

export default useChessStore;
