import useAuthStore from "../store/authStore.ts";

export default function Profile() {

  const user = useAuthStore(state=>state.user);

  if (!user) return <div>Not authenticated</div>;

  return <div>
    {user.username}
  </div>
}
