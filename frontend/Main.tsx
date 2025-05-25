import React from "react";
import { useNavigate } from "react-router-dom";
import { ImageData } from "./App";

interface MainProps {
  images: Map<string, ImageData>;
  user: String;
}

export default function Main({ user, images }: MainProps) {
  const navigate = useNavigate();
  const logout = async () => {
    try {
      console.log("logging out");
      // Send a request to the backend to log the user out
      const response = await fetch('http://localhost:8080/logout', {
        method: 'POST', // or 'DELETE' depending on your backend's method
        credentials: 'include', // Ensure cookies are sent with the request
      });

      if (!response.ok) {
        throw new Error('Failed to log out');
      }

      // Optionally, redirect the user after logging out (e.g., to the login page)
      navigate('/'); // Or use `navigate('/login')` with `useNavigate` in React Router

    } catch (error) {
      console.error('Logout failed:', error);
      // Optionally, handle any errors (e.g., show a message to the user)
    }
  };
  return (
    <div>
      <p>{user}</p>
      <button onClick={() => {navigate('/signup')}}>Singup</button>
      <button onClick={() => {navigate('/login')}}>login</button>
      <button onClick={logout}>logout</button>
      <button onClick={() => {navigate('/admin')}}>Admin</button>
      {Array.from(images).map(([key, val]) => (
        <a onClick={() => {navigate('/single/'+val.FilePath)}}>
          <img src={val.FilePath} alt={val.Description} className="image" />
        </a>
      ))}
    </div>
  );
}
