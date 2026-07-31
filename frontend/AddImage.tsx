import React, { useState } from "react";
import { ImageUploadMetadata } from "./App";
import { useNavigate } from "react-router-dom";

// Define the interface for the form data (image details)


export default function AddImage() {

  const navigate = useNavigate();
  // State to handle form data
  const [formData, setFormData] = useState<ImageUploadMetadata>({
    camera: "",
    aperture: null,
    shutterSpeed: null,
    iso: 0,
    filmStock: null,
    location: "",
    description: ""
  });

  const [file, setFile] = useState<File | null>(null);

  // State to track form submission status
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  // Handle form input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

 
  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    setSuccessMessage(null);

    // Form data validation
    if (
      formData.camera == null ||
      formData.description == null ||
      formData.iso == null ||
      formData.location == null ||
      formData.camera == "" ||
      formData.description == "" ||
      formData.location == "" ||
      !file
    ) {
      setError("Please fill all the fields and select a file.");
      setIsSubmitting(false);
      return;
    }

    const formDataToSend = new FormData();
    formDataToSend.append("file", file);
    formDataToSend.append(
      "metadata",
      new Blob([JSON.stringify(formData)], { type: "application/json" })
    );

    try {
      const response = await fetch("/images", {
        method: "POST",
        body: formDataToSend,
        credentials: "include"
      });

      if (response.ok) {
        setSuccessMessage("Image added successfully!");
        // Clear the form data on success
        setFormData({
          camera: "",
          aperture: null,
          shutterSpeed: null,
          iso: 0,
          filmStock: null,
          location: "",
          description: ""
        });
        // window.location.href = "/";
        window.location.reload();
        // navigate("/");
      } else {
        setError("Failed to add image. Please try again.");
      }
    } catch (err) {
      setError("Error occurred while submitting the form.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="form-container">
      <h2>Add Image</h2>
      {error && <div className="red">{error}</div>}
      {successMessage && <div className="green">{successMessage}</div>}
      <form onSubmit={handleSubmit}>
      <div>
          <label>Camera:</label>
          <input
            className="form-input"
            type="text"
            name="camera"
            value={formData.camera}
            onChange={handleInputChange}
            placeholder="What camera did you use?"
            required
          />
        </div>


        <div>
          <label>ISO:</label>
          <input
            className="form-input"
            type="text"
            name="iso"
            value={formData.iso}
            onChange={handleInputChange}
            placeholder="What ISO?"
            required
          />
        </div>

        <div>
          <label>Shutter Speed:</label>
          <input
            className="form-input"
            type="text"
            name="shutterSpeed"
            value={formData.shutterSpeed || ""}
            onChange={handleInputChange}
            placeholder="What shutter speed?"
          />
        </div>

        <div>
          <label>Aperture:</label>
          <input
            className="form-input"
            type="text"
            name="aperture"
            value={formData.aperture || ""}
            onChange={handleInputChange}
            placeholder="What Aperture?"
          />
        </div>


        <div>
          <label>Film Stock:</label>
          <input
            className="form-input"
            type="text"
            name="filmStock"
            value={formData.filmStock || ""}
            onChange={handleInputChange}
            placeholder="What Film stock?"
          />
        </div>
      
        <div>
          <label>Description:</label>
          <input
            className="form-input"
            type="text"
            name="description"
            value={formData.description}
            onChange={handleInputChange}
            placeholder="What in it?"
            required
          />
        </div>

        <div>
          <label>Location:</label>
          <input
            className="form-input"
            type="text"
            name="location"
            value={formData.location}
            onChange={handleInputChange}
            placeholder="Where'd ya take it?"
            required
          />
        </div>

        
        <div>
          <label>Image File:</label>
          <input
            className="form-input"
            type="file"
            name="file"
            onChange={e => setFile(e.target.files[0])}
            accept="image/*"
            required
          />
        </div>

        <button type="submit" className="small-button" disabled={isSubmitting}>
          {isSubmitting ? "Submitting..." : "Add Image"}
        </button>
      </form>
    </div>
  );
}
