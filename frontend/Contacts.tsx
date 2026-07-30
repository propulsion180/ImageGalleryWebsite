import React, { useEffect, useState } from "react";
import { ImageUploadMetadata } from "./App";
import { useNavigate } from "react-router-dom";


export type ContactRequest = {
    id: number;
    name: string;
    email: string;
    message: string;
    referenceImageId: number | null;
}

export default function Contacts(){
  
  const [contactRequest, setContactRequests] = useState<ContactRequest[] | null>(null);

  function bytesToMegabytes(input: number): number{
    return input / (1024 * 1024);
  }

  useEffect(() => {
    fetch("/contact")
      .then((response) => {
        if(!response.ok){
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json()
      })
      .then((data: ContactRequest[]) => {
        setContactRequests(data);
      })
      .catch((error) => {
        console.error("Error fetching the server stats", error);
      });
  }, []);

  const deleteRequest = async (id: number) => {
    try{
      const response = await fetch(`/contact/${id}`,{
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });

      if(!response.ok){
        throw new Error("Failed to delete the contact request: " + id);
      }

      window.location.reload();
    } catch (error) {
      console.error("Error deleting image: ", error);
      alert("Error deleting image");
    }
  }
  
  return contactRequest !== null ? (
    <>
      <div className="stats-table-wrapper">
        <table className="stats-table">
          <thead>
            <tr><th colSpan={6}>Contact Requests</th></tr>
            <tr><th>ID</th><th>Name</th><th>Email</th><th>Message</th><th>ImageId</th><th>Logout</th></tr>
          </thead>
          <tbody>
            {contactRequest.map((row, i) => (
              <tr key={i}>
                <td>{row.id}</td><td>{row.name}</td><td>{row.email}</td><td>{row.message}</td><td>{row.referenceImageId || ""}</td><td><a className="small-button" onClick={() => {deleteRequest(row.id);}}>Delete</a></td>
              </tr> 
            ))}
          </tbody>
        </table>
      </div>

    </>
  ):(
    <><p><strong>Loading...</strong></p></>    
  );
}
