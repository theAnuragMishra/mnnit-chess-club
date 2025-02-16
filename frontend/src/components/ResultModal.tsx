import useChessStore from "../store/gameStore";

export default function ResultModal({ onClose }: { onClose: () => void }) {
  const result = useChessStore((state) => state.result);
  if (!open) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-lg text-center w-80 border border-gray-300 relative">
        <h2 className="text-2xl font-bold text-gray-800 mb-2">Game Over</h2>
        <p className="text-lg text-gray-600 mb-4">{result}</p>
        <button
          className="absolute top-2 right-2 text-gray-500 hover:text-gray-800"
          onClick={onClose}
        >
          ✖
        </button>
      </div>
    </div>
  );
}
