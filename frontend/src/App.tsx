import { Route, Routes } from 'react-router'
import './App.css'
import Profile from './pages/Profile'
import Login from './pages/Login'
import Signup from './pages/Signup'

function App() {

  return (
    <Routes>
      <Route path="/profile" element=<Profile />></Route>
      <Route path="/login" element=<Login />></Route>
      <Route path="/signup" element=<Signup />></Route>
    </Routes>
  )
}

export default App
