import useAuthStore from "../store/authStore.ts";
import {useNavigate} from "react-router";
import {useEffect} from "react";
import Loading from "../components/Loading.tsx";

export default function Profile() {

  const user = useAuthStore(state=>state.user);
  const loading = useAuthStore(state=>state.loading);
  const navigate = useNavigate();

  useEffect(() => {
    if (!user && !loading) {
      navigate("/login");
    }
  }, [user, navigate, loading]);

  if (loading) {
    return <Loading/>
  }
  

  return <div className="flex-col items-start p-4 text-xl bg-gray-600 rounded-xl m-5 w-3/5">
    <div className="text-3xl">{user?.username}</div>
    <div><p>x Games</p>
      <p>y wins, z losses</p></div>

  </div>
}
