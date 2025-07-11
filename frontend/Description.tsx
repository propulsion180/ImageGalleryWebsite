import React from "react";
import { ImageData } from "./App";

interface DescriptionProps {
  image: ImageData;
}

export default function Description({ image }: DescriptionProps) {
  return (
    <div>
      <table className="description-table">
        {/* <thead>
          <tr>
            <th>Property</th>
            <th>Value</th>
          </tr>
        </thead> */}
        <tbody>
          <tr>
            <td>Description</td>
            <td>{image.description}</td>
          </tr>
          <tr>
            <td>Location</td>
            <td>{image.location}</td>
          </tr>
          <tr>
            <td>Aperture</td>
            <td>f/{image.aperture}</td>
          </tr>
          <tr>
            <td>Shutter Speed</td>
            <td>{image.shutterspeed}s</td>
          </tr>
          <tr>
            <td>ISO</td>
            <td>{image.iso}</td>
          </tr>
          <tr>
            <td>Filepath</td>
            <td>{image.filepath}</td>
          </tr>
        </tbody>
      </table>
    </div>
    // <div>
    //   <p>Description: {image.description}</p>
    //   <p>Location: {image.location}</p>
    //   <p>ShutterSpeed: 1/{image.shutterspeed}s</p>
    //   <p>ISO: {image.iso}</p>
    //   <p>Apeture: f/{image.aperture}</p>
    //   <p>Filepath: {image.filepath}</p>
    // </div>
  );
}
