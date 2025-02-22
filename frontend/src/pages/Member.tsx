import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router";
import { getBaseURL } from "../utils/urlUtils";

export default function Member() {
  const params = useParams();
  const { data } = useSuspenseQuery({
    queryKey: [params.username, "userinfo"],
    queryFn: async () => {
      const response = await fetch(
        `${getBaseURL()}/profile/${params.username}`,
        {
          credentials: "include",
        },
      );
      if (!response.ok) {
        throw new Error("Failed to fetch user data");
      }
      const x = await response.json();
      console.log(x);

      return x; // Convert to JSON
    },
    refetchOnMount: true,
  });

  return (
    <div className="flex-col items-start p-4 text-xl bg-gray-600 rounded-xl m-5 w-3/5">
      <div className="text-3xl mb-2">{params.username}'s Games</div>
      <div>
        {data.map((item) => (
          <Link
            key={item.ID}
            to={`/game/${item.ID}`}
            className="flex gap-2 underline text-blue-500"
          >
            <span>{item.WhiteUsername}</span>{" "}
            <span>{item.Result !== "ongoing" ? item.Result : "*"}</span>
            <span>{item.BlackUsername}</span>
          </Link>
        ))}
      </div>
    </div>
  );
}
