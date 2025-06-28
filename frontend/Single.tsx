import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { ImageData } from "./App";

interface MainProps {
  images: Map<string, ImageData>;
  user: String;
}

export default function Single({ images, user }: MainProps) {
  const { fp } = useParams();

  if (fp == "") {
    const navigate = useNavigate();
    navigate("/");
  }

  return (
    <div>
      <img src={"/" + fp} />
    </div>
  );
}
