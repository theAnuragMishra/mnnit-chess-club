import { useState } from "react";
import useWebSocketStore from "../store/socketStore";
import useAuthStore from "../store/authStore";
import useChessStore from "../store/gameStore";
import useChatStore from "../store/chatStore";

export default function Chat({ gameID }: { gameID: string }) {
  const { sendMessage } = useWebSocketStore();
  const { messages } = useChatStore();
  const { user } = useAuthStore();
  const { whiteID, blackID, whiteusername, blackusername } = useChessStore();
  const [input, setInput] = useState("");

  const handleSendMessage = () => {
    if (input.trim() !== "") {
      sendMessage({
        type: "chat",
        payload: {
          sender: Number(user?.userID),
          receiver: Number(user?.userID === whiteID ? blackID : whiteID),
          senderUsername: user?.username,
          receiverUsername: user?.username === whiteusername ? blackusername : whiteusername,
          text: input,
          gameID: Number(gameID),
        },
      });
      setInput("");
    }
  };

  return (
    <div className="p-4 border rounded shadow-lg w-1/4">
      <div className="h-64 overflow-y-auto border-b mb-2 p-2 text-lg">
        {messages.map((msg, index) => {

          if(msg.gameID == gameID) {return <div key={index} className="mb-1">
            <strong>{msg.sender}:</strong> {msg.text}
          </div>}

        })}
      </div>
      <div className="flex gap-2 w-full text-sm">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          className="border p-2 flex-1 rounded"
          placeholder="Type a message..."
        />
        <button
          onClick={handleSendMessage}
          className="bg-blue-500 text-white px-4 py-2 rounded"
        >
          Send
        </button>
      </div>
    </div>
  );
}
