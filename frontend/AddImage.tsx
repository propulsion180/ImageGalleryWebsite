import React, { useState } from "react";

// Define the interface for the form data (image details)
interface ImageFormData {
  description: string;
  location: string;
  iso: string;
  shutterSpeed: string;
  aperture: string;
  file: File | null;
}

interface AdminProps {
  host: string;
}

export default function AddImage({ host }: AdminProps) {
  // State to handle form data
  const [formData, setFormData] = useState<ImageFormData>({
    description: "",
    location: "",
    iso: "",
    shutterSpeed: "",
    aperture: "",
    file: null,
  });

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

  // Handle file input change
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFormData((prevData) => ({
        ...prevData,
        file: e.target.files[0],
      }));
    }
  };

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    setSuccessMessage(null);

    const { description, location, iso, shutterSpeed, aperture, file } =
      formData;

    // Form data validation
    if (
      !description ||
      !location ||
      !iso ||
      !shutterSpeed ||
      !aperture ||
      !file
    ) {
      setError("Please fill all the fields and select a file.");
      setIsSubmitting(false);
      return;
    }

    const formDataToSend = new FormData();
    formDataToSend.append("description", description);
    formDataToSend.append("location", location);
    formDataToSend.append("iso", iso);
    formDataToSend.append("shutterspeed", shutterSpeed);
    formDataToSend.append("aperture", aperture);
    formDataToSend.append("file", file);

    try {
      const response = await fetch("https://" + host + "/addimage", {
        method: "POST",
        body: formDataToSend,
      });

      if (response.ok) {
        setSuccessMessage("Image added successfully!");
        // Clear the form data on success
        setFormData({
          description: "",
          location: "",
          iso: "",
          shutterSpeed: "",
          aperture: "",
          file: null,
        });
        // window.location.href = "/";
        window.location.reload();
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
          <label>Description:</label>
          <input
            className="form-input"
            type="text"
            name="description"
            value={formData.description}
            onChange={handleInputChange}
            placeholder="Enter image description"
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
            placeholder="Enter location"
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
            placeholder="Enter ISO"
            required
          />
        </div>

        <div>
          <label>Shutter Speed:</label>
          <input
            className="form-input"
            type="text"
            name="shutterSpeed"
            value={formData.shutterSpeed}
            onChange={handleInputChange}
            placeholder="Enter shutter speed"
            required
          />
        </div>

        <div>
          <label>Aperture:</label>
          <input
            className="form-input"
            type="text"
            name="aperture"
            value={formData.aperture}
            onChange={handleInputChange}
            placeholder="Enter aperture"
            required
          />
        </div>

        <div>
          <label>Image File:</label>
          <input
            className="form-input"
            type="file"
            name="file"
            onChange={handleFileChange}
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
