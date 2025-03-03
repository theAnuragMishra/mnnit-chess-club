import useChessStore from "../store/gameStore.ts";
import {useState} from "react";

export default function HistoryTable({history}: {history: any}) {

    const {setGroundFen, setGroundLastMoves, setGroundViewOnly} = useChessStore();
    const [activeIndex, setActiveIndex] = useState(history.length - 1);
    const [active12, setActive12] = useState(history[history.length - 1][1] ? 1 : 0);
    return (
        <div className="overflow-y-auto h-[250px] overflow-x-hidden">
            {history &&
                history.map((move, index) => {
                    return (
                        <div key={index} className="w-full grid grid-cols-[1fr_16fr_16fr] gap-10">
                      <span>
                        {index + 1}
                          {".    "}
                      </span>
                            {move[0] && (
                                <button
                                    onClick={() => {
                                        setGroundFen(move[0].MoveFen);
                                        setGroundLastMoves(move[0].Orig, move[0].Dest);
                                        if (index === history.length - 1 && !move[1]) {
                                            setGroundViewOnly(false);
                                        } else {
                                            setGroundViewOnly(true);
                                        }
                                    }}
                                    className="cursor-pointer"
                                >
                                    {move[0].MoveNotation}
                                </button>
                            )}
                            {move[1] && (
                                <button
                                    onClick={() => {
                                        setGroundFen(move[1].MoveFen);
                                        setGroundLastMoves(move[1].Orig, move[1].Dest);
                                        if (index === history.length - 1) {
                                            setGroundViewOnly(false);
                                        } else {
                                            setGroundViewOnly(true);
                                        }
                                    }}
                                    className="cursor-pointer"
                                >
                                    {move[1].MoveNotation}
                                </button>
                            )}
                        </div>
                    );
                })}
        </div>
    )
}