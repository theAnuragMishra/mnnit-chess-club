import { useState, useEffect, useRef } from "react";

interface ChessTimerProps {
  initialTime: number;
  active: boolean;
  onTimeUp: () => void;
}

export default function Clock({
  initialTime,
  active,
  onTimeUp,
}: ChessTimerProps) {
  // Time is stored in milliseconds.
  const [time, setTime] = useState(initialTime * 1000);
  const intervalRef = useRef<number | null>(null);

  // Reset time when initialTime changes.
  useEffect(() => {
    setTime(initialTime * 1000);
  }, [initialTime]);

  // Start or stop the interval based on the active prop.
  useEffect(() => {
    if (!active) {
      if (intervalRef.current !== null) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      return;
    }

    // Start the timer if active.
    intervalRef.current = window.setInterval(() => {
      setTime((prev) => {
        if (prev <= 100) {
          clearInterval(intervalRef.current!);
          intervalRef.current = null;
          return 0;
        }
        return prev - 100;
      });
    }, 100);

    // Clean up on unmount or when active changes.
    return () => {
      if (intervalRef.current !== null) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [active]);

  // Call onTimeUp when time reaches 0.
  useEffect(() => {
    if (time === 0 && active) onTimeUp();
  }, [time, onTimeUp, active]);

  // Format the time as minutes:seconds.
  const formatTime = (time: number): string => {
    const minutes = Math.floor(time / 60000);
    const seconds = Math.floor((time % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  return (
    <span
      className={`px-2 py-1 my-2 rounded-md text-2xl ${
        active ? "font-bold bg-white text-black" : "bg-black text-gray-400"
      }`}
    >
      {formatTime(time)}
    </span>
  );
}
