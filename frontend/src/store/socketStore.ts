import { create } from "zustand";
import useChessStore from "./gameStore.ts";

// Define WebSocket message type
interface Message {
  type: string;
  payload?: any;
}

// Zustand store type definition
interface WebSocketState {
  connect: () => void;
  setNavigate: (navFunc: any) => void;
  sendMessage: (message: Message) => void;
  close: () => void;
}

// Create Zustand store
const useWebSocketStore = create<WebSocketState>(() => {
  let socket: WebSocket | null = null;
  let navigate: any = null;

  return {
    connect: () => {
      if (socket && socket.readyState === WebSocket.OPEN) return;

      socket = new WebSocket("ws://localhost:8080/ws");

      socket.onopen = () => console.log("Connected to WebSocket");
      socket.onclose = () => console.log("Disconnected from WebSocket");
      socket.onerror = (error) => console.error("WebSocket Error:", error);
      socket.onmessage = (e) => {
        const data = JSON.parse(e.data);

        if (data.type === "Init_Game") {
          navigate(`/game/${data.payload.GameID}`);
        }
        if (data.type === "timeup") {
          console.log(data);
          useChessStore.setState(() => ({ result: data.payload.Result }));
        }
        if (data.type === "Move_Response") {
          // console.log(data);
          const {
            updateGround,
            updateFen,
            updateHistory,
            setTimeBlack,
            setTimeWhite,
          } = useChessStore.getState();
          updateFen(data.payload.move.MoveFen);
          updateGround();
          updateHistory(data.payload.move);
          setTimeWhite(data.payload.timeWhite);
          setTimeBlack(data.payload.timeBlack);
          // console.log(
          //   useChessStore.getState().timeBlack,
          //   useChessStore.getState().timeWhite,
          // );

          if (data.payload.Result !== "") {
            useChessStore.setState(() => ({ result: data.payload.Result }));
          }
        }
      };
    },

    setNavigate: (navFunc: any) => {
      navigate = navFunc;
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
