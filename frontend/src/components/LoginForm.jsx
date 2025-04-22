import { useState } from 'react'

//js login form for the UI
function LoginForm({onLogin}) {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [message, setMessage] = useState('')
    //handle the submit button
    const handleSubmit = async (e) => {
        e.preventDefault()

        //send the json username and pw to Go to process
        const res = await fetch('/login', {
            method: 'POST', 
            headers: {'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })

        //await the response, notify the App if login is successful
        const data = await res.json()
        if (res.ok) {
            setMessage(data.message || 'Login successful.')
            onLogin?.()
        } else {
            setMessage(data.error || 'Login failed.')
        }
        
    }
//return the log in form for the user to input their user name and password ->HTML design (need CSS added later)
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
