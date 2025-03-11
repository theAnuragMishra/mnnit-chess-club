import { Route, Routes } from "react-router";
import "./App.css";
import Member from "./pages/Member.tsx";
import Login from "./pages/Login";
import useAuthStore from "./store/authStore.ts";
import { Suspense, useEffect } from "react";
import RootLayout from "./layouts/RootLayout.tsx";
import Home from "./pages/Home.tsx";
import Loading from "./components/Loading";

import Game from "./pages/Game.tsx";
import Play from "./pages/Play.tsx";
import ServerErrorPage from "./pages/ServerErrorPage.tsx";

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

        <Route
          path="/member/:username"
          element={
            <Suspense fallback={<Loading />}>
              <Member />
            </Suspense>
          }
        />
        <Route path="/login" element={<Login />} />
        <Route path="/play" element={<Play />} />
        <Route
          path="/game/:gameID"
          element={
            <Suspense fallback={<Loading />}>
              <Game />
            </Suspense>
          }
        />
        <Route path="/error-page" element={<ServerErrorPage />} />
      </Route>
    </Routes>
  );
}

export default App;
