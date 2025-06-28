import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { ImageData } from "./App";

interface HeaderProps {
  user: string;
  admin: boolean;
  logout: () => Promise<void>;
}

export default function Header({ user, admin, logout }: HeaderProps) {
  const navigate = useNavigate();

  return (
    <div className="header">
      {user == "" && <h1>Migada's Image Gallery</h1>}
      {user != "" && <h1>Welcome {user}</h1>}
      <div className="navButtonContainer">
        <a
          className="nav-button"
          onClick={() => {
            navigate("/");
          }}
        >
          Home
        </a>
        {user == "" && (
          <a
            className="nav-button"
            onClick={() => {
              navigate("/signup");
            }}
          >
            Signup
          </a>
        )}
        {user == "" && (
          <a
            className="nav-button"
            onClick={() => {
              navigate("/login");
            }}
          >
            Login
          </a>
        )}
        {user != "" && admin && (
          <a
            className="nav-button"
            onClick={() => {
              navigate("/admin");
            }}
          >
            Admin
          </a>
        )}
        {user != "" && (
          <a className="nav-button" onClick={logout}>
            Logout
          </a>
        )}
      </div>
    </div>
  );
}
