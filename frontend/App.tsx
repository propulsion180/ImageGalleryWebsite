import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import ReactDOM from "react-dom";
import {useNavigate, Routes, Route, BrowserRouter as Router } from "react-router-dom";
import Main from "./Main";
import Single from "./Single";
import Login from "./Login";
import Signup from "./Signup";
import Admin from "./Admin";

export type ImageData = {
  FilePath: string;
  Description: string;
  ISO: string;
  ShutterSpeed: string;
  Aperture: string;
  Location: string;
};

const App: React.FC = () => {
  console.log("starting");
  const [images, setImages] = useState<Map<string, ImageData>>(new Map());
  const [user, setUser] = useState<string>("bruh");
  const [admin, setAdmin] = useState<boolean>(false);
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
        data.forEach((image) => newImages.set(image.FilePath, image));
        setImages(newImages);
      })
      .catch((error) => {
        console.error("Error fetching nodes", error);
      });
    console.log("queried");
  }, []);

  

  const host = window.location.host;



  return (
    <div>
      <Router>
        <Routes>
          <Route path="/" element={<Main user={user} images={images}/>} />
          <Route
            path="/single"
            element={<Single images={images} user={user} />}
          />
          <Route path="/login" element={<Login setUser={setUser} />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/admin" element={<Admin admin={admin} images={images}/>} />
        </Routes>
      </Router>
    </div>
  );
};

export default App;
