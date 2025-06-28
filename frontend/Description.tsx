import React from "react";
import { ImageData } from "./App";

interface DescriptionProps {
  image: ImageData;
}

export default function Description({ image }: DescriptionProps) {
  return <p>Filepath: {image.filepath}</p>;
}
