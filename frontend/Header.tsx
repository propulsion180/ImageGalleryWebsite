import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { ImageData, User } from "./App";

interface HeaderProps {
  user: User | null;
  logout: () => Promise<void>;
}

export default function Header({ user, logout }: HeaderProps) {
  const navigate = useNavigate();

  return (
    <div className="header">
      {user == null && <h1>Migada's Image Gallery</h1>}
      {user != null && <h1>Welcome back {user.uName}.</h1>}
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
      <hr className="hr-solid" />
    </div>
  );
}
