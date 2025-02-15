import ChessBoard from "../components/ChessBoard";
import useWebSocketStore from "../store/socketStore";
import { useEffect } from "react";

const GamePage = () => {
  const connect = useWebSocketStore((state) => state.connect);

  useEffect(() => {
    connect();
  }, [connect]);

  return (
    <div className="flex justify-center items-center h-screen">
      <ChessBoard />
    </div>
  );
};

export default GamePage;
