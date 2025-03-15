import { Flag } from "@phosphor-icons/react";
import useWebSocketStore from "../store/socketStore";

export default function DrawResign() {
  const sendMessage = useWebSocketStore((state) => state.sendMessage);
  const handleDraw = () => {
    sendMessage({
      type: "Draw",
      payload: {},
    });
  };
  const handleResign = () => {
    sendMessage({
      type: "Resign",
      payload: {},
    });
  };
  return (
    <div className="text-white w-full flex justify-center items-center gap-2">
      <button
        className="px-4 py-2 rounded-lg hover:bg-gray-600"
        onClick={handleDraw}
      >
        1/2
      </button>
      <button
        onClick={handleResign}
        className="px-4 py-2 rounded-lg hover:bg-gray-600"
      >
        <Flag size={32} />
      </button>
    </div>
  );
}
