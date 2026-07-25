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
import Header from "./Header";

export enum ImageType {
  DIGITAL,
  FILM
}

export type ImageData = {
  id: number;
  filename: string;
  contentType: string;
  fileSizeBytes: number;
  camera: String;
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
  uName: string;
  perms: UserType;
}

const App: React.FC = () => {
  console.log("starting");
  const navigate = useNavigate();
  const [images, setImages] = useState<Map<string, ImageData>>(new Map());
  const [user, setUser] = useState<User | null>(null);
  const logout = async () => {
    try {
      console.log("logging out");
      // Send a request to the backend to log the user out
      const response = await fetch("https://" + host + "/logout", {
        method: "POST",
        credentials: "include", // Ensure cookies are sent with the request
      });

      if (!response.ok) {
        throw new Error("Failed to log out");
      }      
      setUser(null);
      navigate("/");
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };
  // useEffect(() => {
  //   console.log("querying");
  //   fetch("https://" + host + "/all")
  //     .then((response) => {
  //       if (!response.ok) {
  //         throw new Error(`HTTP error! status ${response.status}`);
  //       }
  //       return response.json();
  //     })
  //     .then((data: ImageData[]) => {
  //       console.log(data);
  //       const newImages = new Map<string, ImageData>();
  //       data.forEach((image) => newImages.set(image.filepath, image));
  //       setImages(newImages);
  //     })
  //     .catch((error) => {
  //       console.error("Error fetching nodes", error);
  //     });

    // console.log("queried");
  // }, []);

  useEffect(() => {
    fetch("https://" + host + "/tknlgn")
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
        <Header user={user} admin={admin} logout={logout} />
        <Routes>
          <Route path="/" element={<Main user={user} images={images} />} />
          <Route path="/single" element={<Single />} />
          <Route
            path="/login"
            element={
              <Login setUser={setUser} host={host} />
            }
          />
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
