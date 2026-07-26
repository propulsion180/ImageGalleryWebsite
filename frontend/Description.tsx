import React from "react";
import { ImageData, ImageType } from "./App";

interface DescriptionProps {
  image: ImageData;
}

export default function Description({ image }: DescriptionProps) {
  return (
    <div>
      <table className="description-table">
        <tbody>
          <tr>
            <td>Camera</td>
            <td>{image.camera}</td>
          </tr>
          <tr>
            <td>Description</td>
            <td>{image.description}</td>
          </tr>
          <tr>
            <td>Location</td>
            <td>{image.location}</td>
          </tr>
          {image.type == ImageType.DIGITAL &&
            
            <tr>
              <td>Aperture</td>
              <td>f/{image.aperture}</td>
            </tr>
          }

          {image.type == ImageType.DIGITAL &&
            <tr>
              <td>Shutter Speed</td>
              <td>f/{image.shutterSpeed}</td>
            </tr>
          }
          {image.type == ImageType.FILM &&
            <tr>
              <td>Film Stock</td>
              <td>{image.filmStock}</td>
            </tr>
          }
          <tr>
            <td>ISO</td>
            <td>{image.iso}</td>
          </tr>
          <tr>
            <td>Filepath</td>
            <td>{image.filename}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
