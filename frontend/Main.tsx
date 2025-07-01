import React from "react";
import { useNavigate } from "react-router-dom";
import { ImageData } from "./App";

interface MainProps {
  images: Map<string, ImageData>;
  user: String;
}

export default function Main({ user, images }: MainProps) {
  const navigate = useNavigate();

  console.log(images);

  return (
    <div className="allimage-grid">
      {Array.from(images).map(([fp, image]) => (
        <div className="allimage-item">
          <a
            onClick={() => {
              navigate("/single", { state: { img: image } });
            }}
          >
            <img src={"/" + image.filepath} alt={image.description} />
          </a>
        </div>
      ))}
    </div>
  );
}
