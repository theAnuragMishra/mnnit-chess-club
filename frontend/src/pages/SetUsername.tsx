import {useState} from "react";
import {getBaseURL} from "../utils/urlUtils.ts";

export default function SetUsername() {
    const [username, setUsername] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async () => {
setLoading(true);

    const res = await fetch(`${getBaseURL()}/set-username`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username}),
        credentials: "include"
    });

    const response = await res.json();

    if (!res.ok) {
        setLoading(false);
        setError(response.error);
        return;
    }

    window.location.href="/";

    }

    return <div className="flex flex-col items-center justify-center mt-10 gap-4 p-6 bg-gray-100 rounded-lg shadow-md">
        <input
            type="text"
            placeholder="Enter your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent"
        />

        {error && <p className="text-red-500 text-sm">{error}</p>}

        <button
            className="w-[200px] bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 transition-all disabled:bg-gray-400 disabled:cursor-not-allowed"
            onClick={handleSubmit}
            disabled={loading}
        >
            {loading ? "Submitting..." : "Submit"}
        </button>
    </div>
}