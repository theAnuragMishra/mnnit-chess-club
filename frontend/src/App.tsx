import { Route, Routes } from "react-router";
import "./App.css";
import Profile from "./pages/Profile";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import useAuthStore from "./store/authStore.ts";
import { Suspense, useEffect } from "react";
import RootLayout from "./layouts/RootLayout.tsx";
import Home from "./pages/Home.tsx";
import Loading from "./components/Loading";

import Play from "./pages/Play.tsx";
import Board from "./pages/BoardTest.tsx";

function App() {
  const checkAuth = useAuthStore((state) => state.checkAuth);
  const loading = useAuthStore((state) => state.loading);
  useEffect(() => {
    checkAuth().catch((e) => console.error(e));
  }, [checkAuth]);

  if (loading) {
    return <Loading />;
  }

  return (
    <Routes>
      <Route element={<RootLayout />}>
        <Route path="/" element={<Home />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/board" element={<Board />} />

        <Route
          path="/play/:gameID"
          element={
            <Suspense fallback={<Loading />}>
              <Play />
            </Suspense>
          }
        />
      </Route>
    </Routes>
  );
}

export default App;
