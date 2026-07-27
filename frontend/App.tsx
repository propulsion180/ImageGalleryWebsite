import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import ReactDOM from "react-dom";
import {
  useNavigate,
  Routes,
  Route,
  BrowserRouter as Router,
} from "react-router-dom";
import Main from "./Main";
import Single from "./Single";
import Login from "./Login";
import Header from "./Header";
import Admin from "./Admin";
import Stats from "./Stats";

export enum ImageType {
  DIGITAL,
  FILM
}

export type ImageData = {
  id: number;
  filename: string;
  contentType: string;
  fileSizeBytes: number;
  camera: string;
  type: ImageType;
  aperture: string | null;
  shutterSpeed: string | null;
  iso: number;
  filmStock: string | null;
  location: string;
  description: string;
  uploadedTime: string;
};

export type ImageUploadMetadata = {
  camera: string;
  aperture: string | null;
  shutterSpeed: string | null;
  iso: number;
  filmStock: string | null; 
  location: string;
  description: string; 
}

export enum UserType {
  ADMIN,
  NORMAL,
}

export type User = {
  id: number
  username: string;
  perms: UserType;
}

const App: React.FC = () => {
  console.log("starting");
  const [user, setUser] = useState<User | null>(null);
  const logout = async () => {
    try {
      console.log("logging out");
      // Send a request to the backend to log the user out
      const response = await fetch("/admin/logout", {
        method: "POST",
        credentials: "include", // Ensure cookies are sent with the request
      });

      if (!response.ok) {
        throw new Error("Failed to log out");
      }      
      setUser(null);
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };
    useEffect(() => {
    fetch("http://" + host + "/admin/tknlgn")
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json();
      })
      .then((data: User) => {
        setUser(data);
      })
      .catch((error) => {
        console.error("Error fetching nodes", error);
      });
  }, []);

  const host = window.location.host;

  return (
    <div className="center">
      <Router>
        <Header user={user} logout={logout} />
        <Routes>
          <Route path="/" element={<Main user={user} host={host} />} />
          <Route path="/single" element={<Single host={host} />} />
          <Route
            path="/login"
            element={
              <Login setUser={setUser} host={host} />
            }
          />
          <Route
            path="/admin"
            element={<Admin user={user} host={host} />}
          />
          <Route path="/stats" element={<Stats />} />
        </Routes>
      </Router>
    </div>
  );
};

export default App;
