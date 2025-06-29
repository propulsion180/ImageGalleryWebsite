import React from "react";
import { ImageData } from "./App";

interface DescriptionProps {
  image: ImageData;
}

export default function Description({ image }: DescriptionProps) {
  return (
    <div>
      <p>Description: {image.description}</p>
      <p>Location: {image.location}</p>
      <p>ShutterSpeed: 1/{image.shutterspeed}s</p>
      <p>ISO: {image.iso}</p>
      <p>Apeture: f/{image.aperture}</p>
      <p>Filepath: {image.filepath}</p>
    </div>
  );
}
