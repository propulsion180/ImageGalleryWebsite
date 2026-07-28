import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { ImageData } from "./App";
import Description from "./Description";

interface SingleProps {
  host: string;
}
export default function Single({ host }: SingleProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const imgid: number | undefined = location.state?.id;
  const [image, setImage] = useState<ImageData | null>(null);
  
  useEffect(() => {
    console.log("Retrieve all images for administration");
    fetch("http://" + host + "/images/" + imgid)
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
        <img src={"/images/files/full/" + image.filename} />
        <Description image={image} />
        <a className="small-button" onClick={() => {navigate("/contactform", { state: { id: image.id }});}}>Order Physical Prints</a>
    </div>
    
  ):(
    <div className="single-container"><p><strong>Loading...</strong></p></div>    
  );
}
