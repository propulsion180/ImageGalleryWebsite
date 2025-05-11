import React, { useNavigate } from "react";
import { ImageData } from "./App";

interface MainProps {
  images: Map<string, ImageData>;
}

export default function Main({ images }: MainProps) {
  const navigate = useNavigate();
  return (
    <div>
      {Array.from(images).map(([key, val]) => (
        <a onClick={() => {navigate('/single/'+{val.FilePath})}}>
          <img src={val.FilePath} alt={val.Description} className="image" />
        </a>
      ))}
    </div>
  );
}
