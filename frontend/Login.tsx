import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { User } from "./App";

interface LoginProps {
  setUser: (value: User | null) => void;
  host: string;
}

export default function Login({ setUser, host }: LoginProps) {
  // States to hold form input values
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); // Prevent page reload on form submission

    // Create an object to send as JSON
    const loginData = {
      uName: username,
      plainPassword: password,
    };

    try {
      // Send POST request to the server
      const response = await fetch("https://" + host + "/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(loginData),
      });

      // Handle non-OK responses
      if (!response.ok) {
        const errorMessage = await response.text();
        setError(errorMessage); // Set the error to be displayed
        return;
      }

      const u: User = await response.json();
      alert(`Welcome, ${u.uName}!`); 
      setUser(u);
      navigate("/");
    } catch (err) {
      setError("An error occurred while trying to log in.");
      console.error("Error:", err);
    }
  };

  return (
    <div className="form-container">
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">Username:</label>
          <input
            className="form-input"
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
            className="form-input"
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <p style={{ color: "red" }}>{error}</p>}
        <button type="submit" className="small-button">
          Login
        </button>
      </form>
    </div>
  );
}
