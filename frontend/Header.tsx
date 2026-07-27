import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { ImageData, User, UserType } from "./App";

interface HeaderProps {
  user: User | null;
  logout: () => Promise<void>;
}

export default function Header({ user, logout }: HeaderProps) {
  const navigate = useNavigate();

  async function handleLogout(e: React.MouseEvent<HTMLAnchorElement>){
    e.preventDefault();
    await logout();
    navigate("/");
  }

  console.log(user);
  return (
    <div className="header">
      {user == null && <h1>Migada's Image Gallery</h1>}
      {user != null && <h1>Welcome back {user.username}.</h1>}
      <div className="navButtonContainer">
        <a
          className="nav-button"
          onClick={() => {
            navigate("/");
          }}
        >
          Home
        </a>
        {user != null && user.perms == 'ADMIN' && (
          <a
            className="nav-button"
            onClick={() => {
              navigate("/stats");
            }}
          >
            Stats
          </a>
        )}
        {user != null && user.perms == 'ADMIN' && (
          <a
            className="nav-button"
            onClick={() => {
              navigate("/admin");
            }}
          >
            Admin
          </a>
        )}
        {user != null && (
          <a className="nav-button" onClick={handleLogout}>
            Logout
          </a>
        )}
      </div>
      <hr className="hr-solid" />
    </div>
  );
}
