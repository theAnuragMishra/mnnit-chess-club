import { useState, useEffect } from "react";

interface ChessTimerProps {
  color: "white" | "black";
  initialTime: number;
  active: boolean;
  onTimeUp: (player: "white" | "black") => void;
}

export default function Clock({
  color,
  initialTime,
  active,
  onTimeUp,
}: ChessTimerProps) {
  const [time, setTime] = useState(initialTime * 1000);
  useEffect(() => {
    setTime(initialTime * 1000);
  }, [initialTime]);

  useEffect(() => {
    let timer: NodeJS.Timeout | null = null;

    if (active && time > 0)
      timer = setInterval(() => setTime((prev) => prev - 100), 100);

    return () => {
      if (timer) clearInterval(timer);
    };
  });
  useEffect(() => {
    if (time <= 0) onTimeUp(color);
  }, [time, onTimeUp, color]);
  const formatTime = (time: number): string => {
    const minutes = Math.floor(time / 60000);
    const seconds = Math.floor((time % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  return (
    <span
      className={`px-2 py-1 my-2 rounded-md text-2xl ${active ? "font-bold bg-white  text-black" : "bg-black text-gray-400"}`}
    >
      {formatTime(time)}
    </span>
  );
}
