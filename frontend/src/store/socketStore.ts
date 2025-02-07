import { create } from "zustand";

// Define WebSocket message type
interface WebSocketMessage {
    type: string;
    payload?: {
        "MoveStr":string,
        "GameId":string
    };
}

// Zustand store type definition
interface WebSocketState {
    connect: () => void;
    sendMessage: (message: WebSocketMessage) => void;
    close: () => void;
}

// Create Zustand store
const useWebSocketStore = create<WebSocketState>()(() => {
    let socket: WebSocket | null = null; // WebSocket instance

    return {
        messages: [],

        connect: () => {
            if (socket && socket.readyState === WebSocket.OPEN) return; // Prevent multiple connections

            socket = new WebSocket("ws://localhost:8080/ws");

            socket.onopen = () => console.log("Connected to WebSocket");
            socket.onclose = () => console.log("Disconnected from WebSocket");
            socket.onerror = (error) => console.error("WebSocket Error:", error);


        },

        sendMessage: (message: WebSocketMessage) => {
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
