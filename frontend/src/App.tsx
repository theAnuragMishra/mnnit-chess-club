import { Route, Routes } from 'react-router'
import './App.css'
import Profile from './pages/Profile'
import Login from './pages/Login'
import Signup from './pages/Signup'
import useAuthStore from "./store/authStore.ts";
import {useEffect} from "react";
import RootLayout from "./layouts/RootLayout.tsx";
import Home from "./pages/Home.tsx";
import Loading from "./components/Loading";

function App() {
    const checkAuth = useAuthStore((state) => state.checkAuth);
    const loading = useAuthStore((state) => state.loading);
    useEffect(() => {
        checkAuth().catch(e=>console.error(e));
    }, [checkAuth]);

    if (loading){
        return <Loading />;
    }

  return (
    <Routes>
        <Route element={<RootLayout/>}>
            <Route path="/" element={<Home/>}/>
      <Route path="/profile" element={<Profile/>}/>
      <Route path="/login" element={<Login/>}/>
      <Route path="/signup" element={<Signup/>}/>
        </Route>
    </Routes>
  )
}

export default App
