import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

interface LoginProps {
    setUser: (value: string) => void;
}

export default function Login({setUser}: LoginProps) {
  // States to hold form input values
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); // Prevent page reload on form submission

    // Create an object to send as JSON
    const loginData = {
      username: username,
      password: password,
    };

    try {
      // Send POST request to the server
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(loginData),
      });

      // Handle non-OK responses
      if (!response.ok) {
        const errorMessage = await response.text();
        setError(errorMessage); // Set the error to be displayed
        return;
      }

      // If login is successful, handle the response
      const data = await response.json();
      alert(`Welcome, ${data.username}!`); // Or redirect user after successful login
      setUser(data.username);

      // Optionally, you can store the token or handle redirecting:
      // localStorage.setItem('auth_token', data.auth_token);
      // window.location.href = '/dashboard'; // Redirect to a protected page

    } catch (err) {
      setError('An error occurred while trying to log in.');
      console.error('Error:', err);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">Username:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">Login</button>
      </form>
      <button onClick={() => {navigate('/')}}>Home</button>
    </div>
  );
}