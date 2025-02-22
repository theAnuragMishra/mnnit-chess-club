import { useState, useEffect } from "react";

interface ChessTimerProps {
  initialWhite: number;
  initialBlack: number;
  turn: "white" | "black";
  onTimeUp: (player: "white" | "black") => void;
}

export default function ChessTimer({
  initialWhite,
  initialBlack,
  turn,
  onTimeUp,
}: ChessTimerProps) {
  const [timeWhite, setTimeWhite] = useState(initialWhite);
  const [timeBlack, setTimeBlack] = useState(initialBlack);

  useEffect(() => {
    setTimeWhite(initialWhite);
    setTimeBlack(initialBlack);
  }, [initialWhite, initialBlack]);

  useEffect(() => {
    let timer: NodeJS.Timeout | null = null;

    if (turn === "white" && timeWhite > 0) {
      timer = setInterval(() => setTimeWhite((prev) => prev - 100), 100);
    } else if (turn === "black" && timeBlack > 0) {
      timer = setInterval(() => setTimeBlack((prev) => prev - 100), 100);
    }

    return () => {
      if (timer) clearInterval(timer);
    };
  }, [turn, timeWhite, timeBlack]);

  useEffect(() => {
    if (timeWhite <= 0) onTimeUp("white");
    if (timeBlack <= 0) onTimeUp("black");
  }, [timeWhite, timeBlack, onTimeUp]);

  const formatTime = (time: number): string => {
    const minutes = Math.floor(time / 60000);
    const seconds = Math.floor((time % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  return (
    <div className="flex flex-col items-center space-y-4">
      <div
        className={`p-4 text-2xl ${turn === "white" ? "font-bold text-blue-600" : "text-gray-600"}`}
      >
        ♔ White: {formatTime(timeWhite)}
      </div>
      <div
        className={`p-4 text-2xl ${turn === "black" ? "font-bold text-red-600" : "text-gray-600"}`}
      >
        ♚ Black: {formatTime(timeBlack)}
      </div>
    </div>
  );
}
