import { Link } from "react-router";
import useAuthStore from "../store/authStore.ts";
import { SignOut } from "@phosphor-icons/react";
import { useNavigate } from "react-router";

export default function Nav() {
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();

  async function handleLogout() {
    await logout();
    navigate("/");
  }
  return (
    <div className="flex w-full text-xl px-4 py-4">
      <ul className="flex w-full justify-between">
        <li className="text-2xl">
          <Link to="/">MNNIT Chess Club</Link>
        </li>
        <li className="flex justify-center items-center gap-1 text-2xl">
          {loading ? (
            <div>Loading...</div>
          ) : user ? (
            <>
              <Link to={`/member/${user?.username}`}>{user?.username}</Link>{" "}
              <button onClick={handleLogout}>
                <SignOut className="h-full cursor-pointer" />
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="text-xl text-blue-300 underline">
                Login
              </Link>
            </>
          )}
        </li>
      </ul>
    </div>
  );
}
