export default function Home(){

    function handleInitGame(){
alert("hello")
    }
    return <div className="flex flex-col items-center justify-center h-full ">
        <button onClick={handleInitGame} className="text-4xl bg-gray-500 p-3 rounded-lg cursor-pointer">Play!</button>
    </div>
}