import React from 'react'
import { ImageData } from "./App";


interface DescriptionProps {
    images: Map<string, ImageData>;
}

export default function Description({ images }: DescriptionProps) {
  return (
    <div>Size of images is {images.size}</div>
  )
}
