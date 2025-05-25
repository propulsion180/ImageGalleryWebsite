import React, { useLocation, useNavigate } from "react";
import { useParams } from "react-router-dom";
import { ImageData } from "./App";

interface MainProps {
  images: Map<string, ImageData>;
  user: String;
}

export default function Single({ images, user }: MainProps) {
  const { fp } = useParams();

  if (fp == "") {
    const navigate = useNavigate();
    navigate('/');
  }

  return (
    <div>
      {Array.from(images).map(([key, val]) => (
        <img src={val.FilePath} alt={val.Description} className="image" />
      ))}
    </div>
  );
}
