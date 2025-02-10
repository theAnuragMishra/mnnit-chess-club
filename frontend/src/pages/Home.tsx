import useAuthStore from "../store/authStore.ts";
import Loading from "../components/Loading";
import useWebSocketStore from "../store/socketStore.ts";
import useChessStore from "../store/gameStore.ts";
import {useEffect} from "react";
import {useNavigate} from "react-router";

export default function Home(){
const user = useAuthStore(state => state.user);
const loading = useAuthStore(state => state.loading);

const {sendMessage, connect} = useWebSocketStore();
const {gameID} = useChessStore();
const navigate = useNavigate();

    useEffect(() => {
        if (gameID) {
            navigate(`/play/${gameID}`);
        }
    }, [gameID, navigate]);


    useEffect(() => {
        connect();

    }, [connect]);

    function handleInitGame(){

sendMessage({type:"init_game"})
    }

    if(loading) return <Loading />;

    return <div className="flex flex-col items-center justify-center h-full ">
        {
        user ? <button onClick={handleInitGame} className="text-4xl bg-gray-500 p-3 rounded-lg cursor-pointer">Play!</button> : <p>Login to play</p>
    }

    </div>
}