import { useState, useEffect } from "react";
import useWebSocketStore from "../store/socketStore"; // Import the Zustand store

export default function SocketTest() {
    const [move, setMove] = useState("");
    const [game, setGame] = useState("");
    const { sendMessage, connect, close } = useWebSocketStore();

    useEffect(() => {
        connect(); // Connect WebSocket on mount

        return () => close(); // Disconnect on unmount
    }, [connect, close, sendMessage]);

    const handleSendMessage = () => {
        sendMessage({ type: "move", payload: {MoveStr: move, GameId: game} });
        setMove("");
    };

    return (
        <div>
            <h2>WebSocket Chat</h2>
            From:
            <input
                type="text"
                className="w-full px-3 py-2 border rounded"
                value={move}
                onChange={(e) => setMove(e.target.value)}
            />

            <input
                type="text"
                className="w-full px-3 py-2 border rounded"
                value={game}
                onChange={(e) => setGame(e.target.value)}
            />
            <button onClick={handleSendMessage} className="rounded-md border mt-4 p-4">Send</button>
<button className="rounded-md border mt-4 p-4" onClick={()=>sendMessage({type:"init_game"})}>Start Game</button>

        </div>
    );
}
