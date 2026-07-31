import React, { useState, useEffect } from "react";
import { ImageData, ImageUploadMetadata, User, UserType } from "./App";
import { useNavigate } from "react-router-dom";
import Description from "./Description";
import AddImage from "./AddImage";
import UpdateImage from "./UpdateImage";

interface AdminProps {
  user: User | null
}

export default function Admin({ user }: AdminProps) {
  const navigate = useNavigate();
  const [page, setPage] = useState<string>("desc");
  const [images, setImages] = useState<Map<number, ImageData>>(new Map());
  const [img, setImg] = useState<ImageData | null>(null);
  if (!user || user.perms != 'ADMIN') {
    navigate("/");
  }

  
  useEffect(() => {
    console.log("Retrieve all images for administration");
    fetch("/images/all")
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json();
      })
      .then((data: ImageData[]) => {
        console.log(data);
        const newImages = new Map<number, ImageData>();
        data.forEach((image) => newImages.set(image.id, image));
        setImages(newImages);
      })
      .catch((error) => {
        console.error("Error fetching nodes", error);
      });

   console.log("All images retrieved");
  }, []);


  const deleteImage = async (id: number) => {
    try {
      const response = await fetch("/images/" + id, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to delete image");
      }
      // window.location.href = "/";
      // window.location.reload();
    } catch (error) {
      console.error("Error deleting image: ", error);
      alert("error deleting image");
    }
  };

  return (
    <>
      <div className="admin-buttons-area">
        <a
          className="nav-button"
          onClick={() => {
            setPage("add");
          }}
        >
          Add an Image
        </a>
        <a
          className="nav-button"
          onClick={() => {
            setPage("desc");
          }}
        >
          Description
        </a>
      </div>

      {page == "add" && true && <AddImage />}
      {page == "up" && true && <UpdateImage image={img} />}

      <table className="description-table">
        <thead>
          <tr>
            <th>Filepath:</th>
            <th>Description</th>
            <th>Location</th>
            <th>Delete</th>
            <th>Update</th>
          </tr>
        </thead>
        <tbody>
          {Array.from(images).map(([key, value]) => (
            <tr>
              <td>{value.filename}</td>
              <td>{value.description}</td>
              <td>{value.location}</td>
              <td>
                <a
                  className="small-button"
                  onClick={() => {
                    deleteImage(value.id);
                  }}
                >
                  Delete
                </a>
              </td>
              <td>
                <a
                  className="small-button"
                  onClick={() => {
                    setImg(value);
                    setPage("up");
                  }}
                >
                  Update
                </a>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
