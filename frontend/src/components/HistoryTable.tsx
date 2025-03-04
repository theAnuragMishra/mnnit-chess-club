import useChessStore from "../store/gameStore.ts";
import { useState, useRef, useEffect } from "react";
import { FastForward, Rewind } from "@phosphor-icons/react";

export default function HistoryTable({ history }: { history: any }) {
  const { setGroundFen, setGroundLastMoves, setGroundViewOnly, board } =
    useChessStore();

  const [activeIndex, setActiveIndex] = useState(
    history.length ? history.length - 1 : -1,
  );
  const [active12, setActive12] = useState(
    history.length ? (history[history.length - 1][1] ? 1 : 0) : -1,
  );

  const moveRefs = useRef<(HTMLButtonElement | null)[]>([]); // Store refs for each move

  // Scroll to active move when activeIndex or active12 changes
  useEffect(() => {
    const activeButton = moveRefs.current[activeIndex * 2 + active12]; // Find the right button
    activeButton?.scrollIntoView({ behavior: "smooth", block: "nearest" });
  }, [activeIndex, active12]);

  return (
    <div className="h-[310px] text-lg px-4 py-2  bg-gray-800 relative">
      <div className="overflow-y-auto h-[250px] overflow-x-hidden">
        {history &&
          history.map((move, index) => {
            return (
              <div
                key={index}
                className="w-full grid grid-cols-[1fr_16fr_16fr] gap-10 gap-y-[10px]"
              >
                <span>
                  {index + 1}
                  {".    "}
                </span>
                {move[0] && (
                  <button
                    ref={(el) => (moveRefs.current[index * 2] = el)}
                    onClick={() => {
                      setGroundFen(move[0].MoveFen);
                      setGroundLastMoves(move[0].Orig, move[0].Dest);
                      if (index === history.length - 1 && !move[1]) {
                        setGroundViewOnly(false);
                      } else {
                        setGroundViewOnly(true);
                      }
                      setActiveIndex(index);
                      setActive12(0);
                    }}
                    className={`cursor-pointer ${activeIndex == index && active12 == 0 ? "bg-gray-700" : ""}`}
                  >
                    {move[0].MoveNotation}
                  </button>
                )}
                {move[1] && (
                  <button
                    ref={(el) => (moveRefs.current[index * 2 + 1] = el)}
                    onClick={() => {
                      setGroundFen(move[1].MoveFen);
                      setGroundLastMoves(move[1].Orig, move[1].Dest);
                      if (index === history.length - 1) {
                        setGroundViewOnly(false);
                      } else {
                        setGroundViewOnly(true);
                      }
                      setActiveIndex(index);
                      setActive12(1);
                    }}
                    className={`cursor-pointer ${activeIndex == index && active12 == 1 ? "bg-gray-700" : ""}`}
                  >
                    {move[1].MoveNotation}
                  </button>
                )}
              </div>
            );
          })}
      </div>
      <div className="absolute bottom-2 w-4/5 flex justify-around">
        <button
          className="cursor-pointer w-1/3 flex justify-center items-center hover:bg-gray-700"
          onClick={() => {
            setGroundFen(history[0][0].MoveFen);
            setGroundLastMoves(history[0][0].Orig, history[0][0].Dest);
            setGroundViewOnly(true);
            setActiveIndex(0);
            setActive12(0);
          }}
        >
          <Rewind size={32} />
        </button>
        <button
          className="cursor-pointer w-1/3 flex justify-center items-center hover:bg-gray-700"
          onClick={() => {
            setGroundFen(board.fen());
            if (history[history.length - 1][1]) {
              setGroundLastMoves(
                history[history.length - 1][1].Orig,
                history[history.length - 1][1].Dest,
              );
              setActiveIndex(history.length - 1);
              setActive12(1);
            } else {
              setGroundLastMoves(
                history[history.length - 1][0].Orig,
                history[history.length - 1][0].Dest,
              );
              setActiveIndex(history.length - 1);
              setActive12(0);
            }
            setGroundViewOnly(false);
          }}
        >
          <FastForward size={32} />
        </button>
      </div>
    </div>
  );
}
