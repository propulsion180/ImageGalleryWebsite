import React, { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { ImageData } from "./App";
import Description from "./Description";

export default function Single() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [image, setImage] = useState<ImageData | null>(null);
  
  useEffect(() => {
    console.log("Retrieve all images for administration");
    fetch("/images/" + id)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json();
      })
      .then((data: ImageData) => {
        setImage(data);
      })
      .catch((error) => {
        console.error("Error fetching nodes", error);
      });

   console.log("All images retrieved");
  }, []);

  return image !== null ? (
      
    <div className="single-container"> 
        <div className="protected-image-wrapper"><img onContextMenu={(e) => e.preventDefault()} draggable={false} src={"/images/files/full/" + image.filename} /></div>
        <Description image={image} />
        <a className="small-button" onClick={() => {navigate("/contactform", { state: { id: image.id }});}}>Order Physical Prints</a>
    </div>
    
  ):(
    <div className="single-container"><p><strong>Loading...</strong></p></div>    
  );
}
