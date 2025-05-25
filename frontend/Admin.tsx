import React, { useState } from 'react'
import { ImageData } from "./App";
import { useNavigate } from 'react-router-dom';
import Description from './Description';
import AddImage from './AddImage';
import UpdateImage from './UpdateImage';


interface AdminProps {
  admin: boolean;
  images: Map<string, ImageData>;
}

export default function Admin({ admin, images }: AdminProps) {
  const navigate = useNavigate();
  const [page, setPage] = useState<string>("desc");
  // if(!admin){
  //   navigate('/');
  // }

  return (
    <>
    <h1>The admin page</h1>

    {true && <div>
        {Array.from(images).map(([key, val]) => (
            <div>
              <h2>{val.FilePath}</h2>
              <p>{val.Description}</p>
            </div>
        ))}
        <button onClick={() => {setPage("add")}}>Add Page</button>
        <button onClick={() => {setPage("up")}}>Update Page</button>
        <button onClick={() => {setPage("desc")}}>Description Page</button>
      </div>
      
    }
    { page == "desc" && true && <Description images={images}/>}
    { page == "add" && true && <AddImage />}
    { page == "up" && true && <UpdateImage />}
    </>
  )
}
