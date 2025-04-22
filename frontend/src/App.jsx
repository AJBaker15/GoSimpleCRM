import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import LoginForm from './components/LoginForm'
import Dashboard from './components/Dashboard'

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  return (
    <>
      <div>
        <h1>Coalition CRM</h1>
        {isLoggedIn ? (
          <>
          <Dashboard />
          </>
        ) : (
          <>
          <LoginForm onLogin={() => setIsLoggedIn(true)} />
          </>
        )}
        
        </div>
    </>
  )
}

export default App
