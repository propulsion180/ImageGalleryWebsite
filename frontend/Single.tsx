import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { ImageData } from "./App";
import Description from "./Description";

export default function Single() {
  const location = useLocation();
  const img: ImageData | undefined = location.state?.img;
  if (!img) {
    return <div className="red">Image Not Found</div>;
  }

  return (
    <div className="single-container">
      <img src={"/" + img.filepath} />
      <Description image={img} />
    </div>
  );
}
