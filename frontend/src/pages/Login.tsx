
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";
import useAuthStore from "../store/authStore.ts";
import {useNavigate} from "react-router";

const loginFormSchema = z.object({
  username: z.string().min(1, { message: "Username is required" }),
  password: z.string().min(4, { message: "Password must be at least 4 characters" }),
});

export default function Login() {
  const navigate = useNavigate();
  const login = useAuthStore(state=>state.login)
  // const {login: lg} = useAuthStore();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<z.infer<typeof loginFormSchema>>({
    resolver: zodResolver(loginFormSchema),
  });

  const onSubmit: SubmitHandler<z.infer<typeof loginFormSchema>> = async (data) =>
  {
   try{
     await login(data.username, data.password);
     navigate('/profile');
   } catch (error) {
     //Todo show login error
     console.log(error);
   }
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <h1 className="text-2xl font-bold mb-4">Login</h1>
      <form
        className="bg-white p-6 rounded shadow-md w-full max-w-sm"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="mb-4">
          <label className="block text-gray-700 font-medium mb-2">Username</label>
          <input
            {...register("username")}
            type="text"
            className="w-full px-3 py-2 border rounded"
          />
          {errors.username && (
            <span className="text-red-500 text-sm">{errors.username.message}</span>
          )}
        </div>

        <div className="mb-4">
          <label className="block text-gray-700 font-medium mb-2">Password</label>
          <input
            {...register("password")}
            type="password"
            className="w-full px-3 py-2 border rounded"
          />
          {errors.password && (
            <span className="text-red-500 text-sm">{errors.password.message}</span>
          )}
        </div>

        <button
          type="submit"
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
        >
          Submit
        </button>
      </form>
    </div>
  );
}

