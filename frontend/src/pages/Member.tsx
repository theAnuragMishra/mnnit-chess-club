import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router";
import { getBaseURL } from "../utils/urlUtils";
import { Asterisk } from "@phosphor-icons/react";

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
      // console.log(x);

      return x; // Convert to JSON
    },
    refetchOnMount: true,
  });

  return (
    <div className="flex-col p-4 text-xl bg-black rounded-xl ">
      <div className="text-5xl mb-4 text-center">{params.username}'s Games</div>
      <div className="flex flex-col items-center gap-2 w-full">
        {data &&
          data.map((item, index) => {
            let x = "bg-red-500";
            if (
              (item.WhiteUsername === params.username &&
                item.Result === "1-0") ||
              (item.BlackUsername === params.username && item.Result === "0-1")
            ) {
              x = "bg-green-700";
            } else if (item.Result === "ongoing") {
              x = "bg-gray-600";
            }
            return (
              <Link
                key={item.ID}
                to={`/game/${item.ID}`}
                className="flex w-4/5 gap-2 bg-gray-800 py-4 px-8 rounded-sm"
              >
                <span className="w-1/3 text-left">{item.WhiteUsername}</span>
                <span className={`w-1/3 flex items-center justify-center ${x}`}>
                  {item.Result !== "ongoing" ? item.Result : <Asterisk />}
                </span>
                <span className="w-1/3 text-right">{item.BlackUsername}</span>
              </Link>
            );
          })}
      </div>
    </div>
  );
}
