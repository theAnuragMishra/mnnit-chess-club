import { create } from "zustand";
import useChessStore from "./gameStore.ts";

interface Move {
  MoveStr: string;
  GameID: string;
}

// Define WebSocket message type
interface Message {
  type: string;
  payload?: any;
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
          console.log(data);
          useChessStore.setState((state) => ({
            gameID: data.payload.GameID,
            player1username: data.payload.player1username,
            player2username: data.payload.player2username,
          }));
        }
        if (data.type === "Result_Alert") {
          console.log(data);
          if (data.payload.Result === "1-0")
            useChessStore.setState((state) => ({ result: "1-0" }));
          else if (data.payload.Result === "0-1")
            useChessStore.setState((state) => ({ result: "0-1" }));
          else if (data.payload.Result === "1/2-1/2")
            useChessStore.setState((state) => ({ result: "1/2-1/2" }));
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
