import React, { useState } from "react";
interface SignupProps {
  host: string;
}
export default function Signup({ host }: SignupProps) {
  // States to hold form input values
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); // Prevent page reload on form submission

    // Create an object to send as JSON
    const signupData = {
      username: username,
      password: password,
    };

    try {
      // Send POST request to the server (Signup endpoint)
      const response = await fetch("http://" + host + "/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(signupData),
      });

      // Handle non-OK responses
      if (!response.ok) {
        const errorMessage = await response.text();
        setError(errorMessage); // Set the error to be displayed
        setSuccessMessage(""); // Clear success message on error
        return;
      }

      // If signup is successful, display success message
      setSuccessMessage("Signup successful! Please log in.");
      setError(""); // Clear any existing error messages
    } catch (err) {
      setError("An error occurred while trying to sign up.");
      setSuccessMessage(""); // Clear success message on error
      console.error("Error:", err);
    }
  };

  return (
    <div className="form-container">
      <h2>Signup</h2>
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
        {successMessage && <p style={{ color: "green" }}>{successMessage}</p>}
        <button type="submit" className="small-button">
          Sign Up
        </button>
      </form>
    </div>
  );
}
