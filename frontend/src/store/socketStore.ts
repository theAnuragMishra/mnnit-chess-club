import { create } from "zustand";
import useChessStore from "./gameStore.ts";
import useAuthStore from "./authStore.ts";

interface Move {
    MoveStr: string;
    GameID: string
}

// Define WebSocket message type
interface Message {
    type: string;
    payload?: Move;
}

// Zustand store type definition
interface WebSocketState {
    connect: () => void;
    sendMessage: (message: Message) => void;
    close: () => void;
}

// Create Zustand store
const useWebSocketStore = create<WebSocketState>((set, get) => {
    let socket: WebSocket | null = null;

    return {
        connect: () => {
            if (socket && socket.readyState === WebSocket.OPEN) return;

            socket = new WebSocket("ws://localhost:8080/ws");

            socket.onopen = () => console.log("Connected to WebSocket");
            socket.onclose = () => console.log("Disconnected from WebSocket");
            socket.onerror = (error) => console.error("WebSocket Error:", error);
            socket.onmessage = (e) => {
                const data = JSON.parse(e.data);
                if (data.type === "move") {
                    useChessStore.getState().makeMove(data.move);
                }
                if (data.type === "Init_Game") {
                    console.log("Game ID:", data);
                    useChessStore.setState((state) => ({
                        gameID: data.payload.GameID,
                        player1: data.payload.player1,
                        player2: data.payload.player2,
                    }));
                }
            };
        },

        sendMessage: (message: Message) => {
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify(message));
            } else {
                console.warn("WebSocket is not connected yet.");
            }
        },

        close: () => {
            if (socket) {
                socket.close();
                socket = null;
            }
        },
    };
});


export default useWebSocketStore;
