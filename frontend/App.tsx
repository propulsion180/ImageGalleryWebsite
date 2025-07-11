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
import Signup from "./Signup";
import Admin from "./Admin";
import Header from "./Header";

export type ImageData = {
  filepath: string;
  description: string;
  iso: string;
  shutterspeed: string;
  aperture: string;
  location: string;
};

interface LoginResponse {
  message: string;
  username: string;
  admin: boolean;
}

const App: React.FC = () => {
  console.log("starting");
  const [images, setImages] = useState<Map<string, ImageData>>(new Map());
  const [user, setUser] = useState<string>("");
  const [admin, setAdmin] = useState<boolean>(false);
  const logout = async () => {
    try {
      console.log("logging out");
      // Send a request to the backend to log the user out
      const response = await fetch("http://" + host + "/logout", {
        method: "POST", // or 'DELETE' depending on your backend's method
        credentials: "include", // Ensure cookies are sent with the request
      });

      if (!response.ok) {
        throw new Error("Failed to log out");
      }
      setUser("");
      setAdmin(false);
    } catch (error) {
      console.error("Logout failed:", error);
      // Optionally, handle any errors (e.g., show a message to the user)
    }
  };
  useEffect(() => {
    console.log("querying");
    fetch("http://" + host + "/all")
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json();
      })
      .then((data: ImageData[]) => {
        console.log(data);
        const newImages = new Map<string, ImageData>();
        data.forEach((image) => newImages.set(image.filepath, image));
        setImages(newImages);
      })
      .catch((error) => {
        console.error("Error fetching nodes", error);
      });

    console.log("queried");
  }, []);

  useEffect(() => {
    fetch("http://" + host + "/tknlgn")
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json();
      })
      .then((data: LoginResponse) => {
        console.log(data.message);
        setUser(data.username);
        setAdmin(data.admin);
      })
      .catch((error) => {
        console.error("Error fetching nodes", error);
      });
  }, []);

  const host = window.location.host;

  return (
    <div className="center">
      <Router>
        <Header user={user} admin={admin} logout={logout} />
        <Routes>
          <Route path="/" element={<Main user={user} images={images} />} />
          <Route path="/single" element={<Single />} />
          <Route
            path="/login"
            element={
              <Login setUser={setUser} setAdmin={setAdmin} host={host} />
            }
          />
          <Route path="/signup" element={<Signup host={host} />} />
          <Route
            path="/admin"
            element={<Admin admin={admin} images={images} host={host} />}
          />
        </Routes>
      </Router>
    </div>
  );
};

export default App;
