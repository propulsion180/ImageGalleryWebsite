import React, { useEffect, useState, useNavigate, useLocation } from "react";
import ReactDOM from "react-dom";
import { Routes, Route, BrowserRouter as Router } from "react-router-dom";
import Main from "./Main";
import Single from "./Single";

export type ImageData = {
  FilePath: string;
  De;
  scription: string;
  ISO: string;
  ShutterSpeed: string;
  Aperture: string;
  Location: string;
};

const App: React.FC = () => {
  console.log("starting");
  const navigate = useNavigate();
  const location = useLocation();
  const [images, setImages] = useState<Map<string, ImageData>>(new Map());
  const [user, setUser] = useState<string>("");
  const [admin, setAdmin] = useState<Boolean>(false);
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
          <Route path="/" element={<Main images={images} />} />
          <Route
            path="/single"
            element={<Single images={images} user={user} />}
          />
        </Routes>
      </Router>
    </div>
  );
};

export default App;
