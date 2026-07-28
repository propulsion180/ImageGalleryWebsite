import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { User } from "./App";
import { ContactRequest } from "./Contacts";


export default function ContactForm() {
  const location = useLocation();
  const contactid: number | null = location.state?.id;
  // States to hold form input values

  const [contactRequest, setContactRequest] = useState<ContactRequest>({
    id: -1,
    name: "",
    email: "",
    message: "",
    referenceImageId: contactid,
  });
  
  const [error, setError] = useState("");
  const navigate = useNavigate();


  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setContactRequest((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };
  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); // Prevent page reload on form submission

    // Create an object to send as JSON
    

    try {
      // Send POST request to the server
      const response = await fetch(`/contact`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(contactRequest),
      });

      // Handle non-OK responses
      if (!response.ok) {
        const errorMessage = await response.text();
        setError(errorMessage); // Set the error to be displayed
        return;
      }

      navigate("/");
    } catch (err) {
      setError("An error occurred while trying to log in.");
      console.error("Error:", err);
    }
  };

  return (
    <div className="form-container">
      <h2>Contact Me</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="name">Name:</label>
          <input
            className="form-input"
            type="text"
            id="name"
            name="name"
            value={contactRequest.name}
            onChange={handleInputChange}
            required
          />
        </div>
        
        <div>
          <label htmlFor="email">Email:</label>
          <input
            className="form-input"
            type="text"
            id="email"
            name="email"
            value={contactRequest.email}
            onChange={handleInputChange}
            required
          />
        </div>
        
        <div>
          <label htmlFor="message">Message:</label>
          <input
            className="form-input"
            type="text"
            id="message"
            name="message"
            value={contactRequest.message}
            onChange={handleInputChange}
            required
          />
        </div>
        {error && <p style={{ color: "red" }}>{error}</p>}
        <button type="submit" className="small-button">
          Send!
        </button>
      </form>
    </div>
  );
}
