import React, { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { ImageData } from "./App";

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
  host: string;
}

export default function UpdateImage({ image, host }: UpdateImageProps) {
  const location = useLocation();
  const navigate = useNavigate();
  // State for form fields

  const [desc, setDesc] = useState<string>(image.description);
  const [loc, setLoc] = useState<string>(image.location);
  const [iso, setIso] = useState<string>(image.iso);
  const [ss, setSs] = useState<string>(image.shutterspeed);
  const [apt, setApt] = useState<string>(image.aperture);

  // State to track form submission status
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
    if (!desc || !loc || !iso || !ss || !apt) {
      setError("Please fill all the fields and select a file.");
      setIsSubmitting(false);
      return;
    }

    const formDataToSend = new FormData();
    formDataToSend.append("description", desc);
    formDataToSend.append("location", loc);
    formDataToSend.append("iso", iso);
    formDataToSend.append("shutterspeed", ss);
    formDataToSend.append("aperture", apt);
    formDataToSend.append("filepath", image.filepath);

    try {
      const response = await fetch("https://" + host + "/setimage", {
        method: "PUT",
        body: formDataToSend,
      });

      if (response.ok) {
        setSuccessMessage("Image updated successfully!");
        navigate("/admin"); // Redirect to another page after success
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
          <label>Description:</label>
          <input
            className="form-input"
            type="text"
            name="description"
            value={desc}
            onChange={(e) => {
              setDesc(e.target.value);
            }}
            required
          />
        </div>

        <div>
          <label>Location:</label>
          <input
            className="form-input"
            type="text"
            name="location"
            value={loc}
            onChange={(e) => {
              setLoc(e.target.value);
            }}
            required
          />
        </div>

        <div>
          <label>ISO:</label>
          <input
            className="form-input"
            type="text"
            name="iso"
            value={iso}
            onChange={(e) => {
              setLoc(e.target.value);
            }}
            required
          />
        </div>

        <div>
          <label>Shutter Speed:</label>
          <input
            className="form-input"
            type="text"
            name="shutterSpeed"
            value={ss}
            onChange={(e) => {
              setSs(e.target.value);
            }}
            required
          />
        </div>

        <div>
          <label>Aperture:</label>
          <input
            className="form-input"
            type="text"
            name="aperture"
            value={apt}
            onChange={(e) => {
              setApt(e.target.value);
            }}
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
