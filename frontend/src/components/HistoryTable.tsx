import useChessStore from "../store/gameStore.ts";
import { useRef, useEffect, Fragment } from "react";
import { FastForward, Rewind } from "@phosphor-icons/react";

export default function HistoryTable({ history }: { history: any }) {
  const {
    setGroundFen,
    setGroundLastMoves,
    setGroundViewOnly,
    board,
    activeIndex,
    setActiveIndex,
  } = useChessStore();
  const moveRefs = useRef<(HTMLButtonElement | null)[]>([]); // Store refs for each move

  // Scroll to active move when activeIndex or active12 changes
  useEffect(() => {
    const activeButton = moveRefs.current[activeIndex]; // Find the right button
    activeButton?.scrollIntoView({ behavior: "smooth", block: "nearest" });
  }, [activeIndex]);

  return (
    <div className="h-[310px] text-lg px-4 py-2  bg-gray-800 relative">
      <div className="overflow-y-auto h-[250px] overflow-x-hidden w-full grid grid-cols-[1fr_16fr_16fr] content-start place-items-center">
        {history &&
          history.map((move, index) => {
            let x;
            if (index % 2 == 0) x = <span>{index / 2 + 1}.</span>;
            return (
              <Fragment key={index}>
                {x}
                <button
                  ref={(el) => (moveRefs.current[index] = el)}
                  onClick={() => {
                    setGroundFen(move.MoveFen);
                    setGroundLastMoves(move.Orig, move.Dest);
                    if (index === history.length - 1) {
                      setGroundViewOnly(false);
                    } else {
                      setGroundViewOnly(true);
                    }
                    setActiveIndex(index);
                  }}
                  className={`cursor-pointer h-fit w-fit px-4 py-2 ${activeIndex == index ? "bg-gray-700" : ""}`}
                >
                  {move.MoveNotation}
                </button>
              </Fragment>
            );
          })}
      </div>
      {history && (
        <div className="absolute bottom-2 w-4/5 flex justify-around">
          <button
            className="cursor-pointer w-1/3 flex justify-center items-center hover:bg-gray-700"
            onClick={() => {
              setGroundFen(history[0].MoveFen);
              setGroundLastMoves(history[0].Orig, history[0].Dest);
              setGroundViewOnly(true);
              setActiveIndex(0);
            }}
          >
            <Rewind size={32} />
          </button>
          <button
            className="cursor-pointer w-1/3 flex justify-center items-center hover:bg-gray-700"
            onClick={() => {
              setGroundFen(board.fen());
              setGroundLastMoves(
                history[history.length - 1].Orig,
                history[history.length - 1].Dest,
              );
              setActiveIndex(history.length - 1);

              setGroundViewOnly(false);
            }}
          >
            <FastForward size={32} />
          </button>
        </div>
      )}
    </div>
  );
}
