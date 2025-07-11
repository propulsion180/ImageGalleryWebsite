import React, { useState } from "react";
import { ImageData } from "./App";
import { useNavigate } from "react-router-dom";
import Description from "./Description";
import AddImage from "./AddImage";
import UpdateImage from "./UpdateImage";

interface AdminProps {
  admin: boolean;
  images: Map<string, ImageData>;
}

export default function Admin({ admin, images }: AdminProps) {
  const navigate = useNavigate();
  const [page, setPage] = useState<string>("desc");
  const [img, setImg] = useState<ImageData>({
    filepath: "",
    description: "",
    location: "",
    iso: "",
    shutterspeed: "",
    aperture: "",
  });
  if (!admin) {
    navigate("/");
  }

  const deleteImage = async (filepath: string) => {
    try {
      const response = await fetch("http://localhost:8080/delimage", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ filepath }),
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to delete image");
      }
      // window.location.href = "/";
      window.location.reload();
    } catch (error) {
      console.error("Error deleting image: ", error);
      alert("error deleting image");
    }
  };

  return (
    <>
      <h1>The admin page</h1>
      {page == "add" && true && <AddImage />}
      {page == "up" && true && <UpdateImage image={img} />}
      {true && (
        <div>
          {Array.from(images).map(([key, val]) => (
            <div>
              <Description image={val} />
              <p>{val.description}</p>
              <button
                onClick={() => {
                  setImg(val);
                  setPage("up");
                }}
              >
                Update
              </button>
              <button
                onClick={() => {
                  deleteImage(val.filepath);
                }}
              >
                Delete
              </button>
            </div>
          ))}
          <button
            onClick={() => {
              setPage("add");
            }}
          >
            Add Page
          </button>

          <button
            onClick={() => {
              setPage("desc");
            }}
          >
            Description Page
          </button>
        </div>
      )}
    </>
  );
}
