import React, { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { ImageData, ImageUploadMetadata } from "./App";

// Define the interface for the form data (image details)
interface ImageFormData {
  description: string;
  location: string;
  iso: string;
  shutterSpeed: string;
  aperture: string;
}

interface UpdateImageProps {
  image: ImageData;
}

export default function UpdateImage({ image }: UpdateImageProps) {
  const location = useLocation();
  const navigate = useNavigate();
  // State for form fields
  const [formData, setFormData] = useState<ImageUploadMetadata>({
    camera: image.camera,
    aperture: image.aperture,
    shutterSpeed: image.shutterSpeed,
    iso: image.iso,
    filmStock: image.filmStock,
    location: image.location,
    description: image.description
  });

  
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };  

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

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
      formData.location == "" 
    ) {
      setError("Please fill all the fields and select a file.");
      setIsSubmitting(false);
      return;
    }
    
    try {
      const response = await fetch("/images/" + image.id, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
        credentials: "include"
      });

      if (response.ok) {
        setSuccessMessage("Image updated successfully!");
        window.location.reload();
        // navigate("/admin"); // Redirect to another page after success
      } else {
        setError("Failed to update image. Please try again.");
      }
    } catch (err) {
      setError("Error occurred while submitting the form.");
    } finally {
      setIsSubmitting(false);
    }
  };

  useEffect(() => {
    if (!image) {
      setError("No image data available");
    }
  }, [image]);

  return (
    <div className="form-container">
      <h2>Update Image</h2>
      {error && <div>{error}</div>}
      {successMessage && <div>{successMessage}</div>}
      <form onSubmit={handleSubmit}>

        <div>
          <label>Camera:</label>
          <input
            className="form-input"
            type="text"
            name="camera"
            value={formData.camera}
            onChange={handleInputChange}
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
            required
          />
        </div>

        

        <button type="submit" className="small-button" disabled={isSubmitting}>
          {isSubmitting ? "Submitting..." : "Update Image"}
        </button>
      </form>
    </div>
  );
}
