import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";
import useAuthStore from "../store/authStore.ts";
import { Link, useNavigate, useSearchParams } from "react-router";
import Loading from "../components/Loading";
import { useEffect } from "react";

const loginFormSchema = z.object({
  username: z.string().min(1, { message: "Username is required" }),
  password: z
    .string()
    .min(4, { message: "Password must be at least 4 characters" }),
});

export default function Login() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const login = useAuthStore((state) => state.login);

  const [searchParams] = useSearchParams();
  const message = searchParams.get("message");

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<z.infer<typeof loginFormSchema>>({
    resolver: zodResolver(loginFormSchema),
  });

  useEffect(() => {
    if (user && !loading) {
      navigate(`/member/${user?.username}`);
    }
  }, [user, navigate, loading]);

  if (loading) return <Loading />;

  const onSubmit: SubmitHandler<z.infer<typeof loginFormSchema>> = async (
    data,
  ) => {
    try {
      await login(data.username, data.password);
      navigate(`/member/${user?.username}`);
    } catch (error) {
      //Todo show login error
      console.log(error);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-full mt-30">
      <h1 className="text-2xl font-bold mb-2">Login</h1>
      {message == "registered" && (
        <p className="mb-2 text-xl text-green-600">
          Account created successfully, now log in!
        </p>
      )}
      <form
        className="bg-gray-800 p-6 rounded shadow-md w-full max-w-sm"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="mb-4">
          <label className="block font-medium mb-2">Username</label>
          <input
            {...register("username")}
            type="text"
            className="w-full px-3 py-2 border rounded"
          />
          {errors.username && (
            <span className="text-red-500 text-sm">
              {errors.username.message}
            </span>
          )}
        </div>

        <div className="mb-4">
          <label className="block font-medium mb-2">Password</label>
          <input
            {...register("password")}
            type="password"
            className="w-full px-3 py-2 border rounded"
          />
          {errors.password && (
            <span className="text-red-500 text-sm">
              {errors.password.message}
            </span>
          )}
        </div>

        <button
          type="submit"
          className="w-full bg-gray-500 text-white py-2 rounded hover:bg-gray-600 cursor-pointer"
        >
          Submit
        </button>
        <Link
          className="w-full text-center block mt-2 text-blue-400 hover:underline"
          to={"/signup"}
        >
          New here? Signup!
        </Link>
      </form>
    </div>
  );
}
