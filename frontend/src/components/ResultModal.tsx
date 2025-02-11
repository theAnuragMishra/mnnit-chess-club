import useGameStore from "../store/gameStore";

const ResultModal = () => {
    const { result } = useGameStore();
    if (!result) return null;

    return (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
            <div className="bg-white p-6 rounded-lg shadow-lg">
                <h2 className="text-2xl text-black">{result}</h2>
            </div>
        </div>
    );
};

export default ResultModal;
