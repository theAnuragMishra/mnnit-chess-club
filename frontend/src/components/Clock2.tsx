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
  const [time, setTime] = useState(initialTime * 1000);
  const animationFrameRef = useRef<number | null>(null);
  const startTimeRef = useRef<number | null>(null);

  useEffect(() => {
    setTime(initialTime * 1000);
  }, [initialTime]);

  useEffect(() => {
    if (active) {
      startTimeRef.current = performance.now();
      const tick = (currentTime: number) => {
        if (!startTimeRef.current) return;
        const elapsed = currentTime - startTimeRef.current;
        const newTime = initialTime * 1000 - elapsed;

        if (newTime <= 0) {
          setTime(0);
          onTimeUp();
          return;
        }

        setTime(newTime);
        animationFrameRef.current = requestAnimationFrame(tick);
      };

      animationFrameRef.current = requestAnimationFrame(tick);
    } else {
      if (animationFrameRef.current !== null) {
        cancelAnimationFrame(animationFrameRef.current);
        animationFrameRef.current = null;
        startTimeRef.current = null;
      }
    }

    return () => {
      if (animationFrameRef.current !== null) {
        cancelAnimationFrame(animationFrameRef.current);
        animationFrameRef.current = null;
        startTimeRef.current = null;
      }
    };
  }, [active, initialTime, onTimeUp]);

  useEffect(() => {
    if (time === 0 && active) onTimeUp();
  }, [time, onTimeUp, active]);

  const formatTime = (time: number): string => {
    const minutes = Math.floor(time / 60000);
    const seconds = Math.floor((time % 60000) / 1000);
    const milliseconds = Math.floor((time % 1000) / 10);
    if (time > 10000)
      return `${minutes}:${seconds.toString().padStart(2, "0")}`;
    return `${minutes}:${seconds.toString().padStart(2, "0")}.${milliseconds.toString().padStart(2, "0")}`;
  };

  return (
    <span
      className={`px-2 py-1 my-2 rounded-md text-4xl ${
        active ? "font-bold " : "bg-black text-gray-400"
      } ${active && time > 10000 ? "bg-white text-black" : ""} ${time < 10000 ? "bg-red-500 text-white" : ""}`}
    >
      {formatTime(time)}
    </span>
  );
}
