import { useState } from 'react'

function LoginForm() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [message, setMessage] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault()

        const res = await fetch('/login', {
            method: 'POST', 
            headers: {'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })

        const data = await res.json()
        if (res.ok) {
            setMessage(data.message || 'Login successful.')
        } else {
            setMessage(data.error || 'Login failed.')
        }
        
    }

    return (
        <div>
            <h2>Log In</h2>
            <form onSubmit={handleSubmit}>
                <label>
                    Username:
                    <input type="text" value={username} onChange={e => setUsername(e.target.value)} />
                </label>
                <br />
                <label>
                    Password:
                    <input type="text" value={password} onChange={e => setPassword(e.target.value)} />
                </label>
                <br />
                <button type="submit">Log In</button>
            </form>
            {message && <p>{message}</p>}
        </div>
    )
}

export default LoginForm
