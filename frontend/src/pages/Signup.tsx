import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";
import useAuthStore from "../store/authStore.ts";
import { Link, useNavigate } from "react-router";
import Loading from "../components/Loading.tsx";
import { useEffect } from "react";

// Define the schema for the signup form
const signupFormSchema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
  username: z
    .string()
    .min(3, { message: "Username must be at least 3 characters" }),
  password: z
    .string()
    .min(6, { message: "Password must be at least 6 characters" }),
});

type SignupFormInputs = z.infer<typeof signupFormSchema>;

export default function Signup() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const loading = useAuthStore((state) => state.loading);
  const signup = useAuthStore((state) => state.register);
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignupFormInputs>({
    resolver: zodResolver(signupFormSchema),
  });

  useEffect(() => {
    if (user && !loading) {
      navigate(`/member/${user?.username}`);
    }
  }, [user, navigate, loading]);

  if (loading) return <Loading />;

  const onSubmit: SubmitHandler<SignupFormInputs> = async (data) => {
    //console.log("Signup Data:", data);
    try {
      await signup(data.username, data.password);
      navigate("/login?message=registered");
    } catch (err) {
      //Todo show registration error
      console.log(err);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center mt-20">
      <h1 className="text-2xl font-bold mb-4">Signup</h1>
      <form
        className="bg-gray-800 p-6 rounded shadow-md w-full max-w-sm"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="mb-4">
          <label className="block font-medium mb-2">Email</label>
          <input
            {...register("email")}
            type="email"
            className="w-full px-3 py-2 border rounded"
          />
          {errors.email && (
            <span className="text-red-500 text-sm">{errors.email.message}</span>
          )}
        </div>

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
          Signup
        </button>
        <Link
          className="w-full text-center block mt-2 text-blue-400 hover:underline"
          to={"/login"}
        >
          Already registered? Login!
        </Link>
      </form>
    </div>
  );
}
