import React, { useEffect, useState } from "react";
import { ImageUploadMetadata } from "./App";
import { useNavigate } from "react-router-dom";

export type LoginAttemptDTO = {
  username: string;
  timestamp: string;
}

export type MonitoringDTO = {
  diskUsedBytes: number;
  diskTotalBytes: number;
  jvmUsedMemoryBytes: number;
  jvmMaxMemoryBytes: number;
  imageCount: number;
  totalImageStorageBytes: number;
  recentErrors: string[];
  recentFailedLogins: LoginAttemptDTO[];
  lastSuccessfulLogin: string;
  requestCountLastHour: number;
  errorCountLastHour: number;
}

export default function Stats(){
  
  const [monitorData, setMonitorData] = useState<MonitoringDTO | null>(null);

  function bytesToMegabytes(input: number): number{
    return input / (1024 * 1024);
  }

  useEffect(() => {
    fetch("/admin/metrics/stats")
      .then((response) => {
        if(!response.ok){
          throw new Error(`HTTP error! status ${response.status}`);
        }
        return response.json()
      })
      .then((data: MonitoringDTO) => {
        setMonitorData(data);
      })
      .catch((error) => {
        console.error("Error fetching the server stats", error);
      });
  }, []);
  
  return (
    <>

      <table className="description-table">
        <thead>
          <tr>
            <th colSpan={2}> System </th>
          </tr>
          <tr>
            <th>System</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          <tr key="Disk usage">
            <td>Disk space used</td>
            <td>{bytesToMegabytes(monitorData?.diskUsedBytes || 0)}MB / {bytesToMegabytes(monitorData?.diskTotalBytes || 0)}MB used</td>
          </tr>

          <tr key="jvmmem">
            <td>JVM Memory used</td>
            <td>{bytesToMegabytes(monitorData?.jvmUsedMemoryBytes || 0)}MB / {bytesToMegabytes(monitorData?.jvmMaxMemoryBytes || 0)}MB used</td>
          </tr>
        </tbody>
      </table>


      
      <table className="description-table">
        <thead>
          <tr>
            <th colSpan={2}> Images </th>
          </tr>
          <tr>
            <th>Metric</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          <tr key="numimages">
            <td>Number of Images</td>
            <td>{monitorData?.imageCount} images</td>
          </tr>

          <tr key="sizimage">
            <td>Size of images</td>
            <td>{bytesToMegabytes(monitorData?.totalImageStorageBytes || 0)}MB of images</td>
          </tr>
        </tbody>
      </table>

      
      <table className="description-table">
        <thead>
          <tr>
            <th colSpan={2}> Request Metrics </th>
          </tr>
          <tr>
            <th>Metric</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          <tr key="lhreq">
            <td>Requests in the last hour</td>
            <td>{monitorData?.requestCountLastHour}</td>
          </tr>

          <tr key="lherr">
            <td>Errors in last hour</td>
            <td>{monitorData?.errorCountLastHour}</td>
          </tr>
        </tbody>
      </table>

      <div className="stats-table-wrapper">
        <table className="stats-table">
          <thead>
            <tr><th colSpan={2}>Login Request Attempts</th></tr>
            <tr>
              <th>Username</th>
              <th>TimeStamp</th>
            </tr>
          </thead>
          <tbody>
            {monitorData?.recentFailedLogins.map((row, i) => (
              <tr key={i}>
                <td>{row.username}</td>
                <td>{new Date(row.timestamp).toLocaleString()}</td>
              </tr> 
            ))}
          </tbody>
        </table>
      </div>
       
      <div className="stats-table-wrapper">
        <table className="stats-table">
          <thead>
            <tr><th>Errors</th></tr>
          </thead>
          <tbody>
            {monitorData?.recentErrors.map((row, i) => (
              <tr key={i}>
                <td>{row}</td>
              </tr> 
            ))}
          </tbody>
        </table>
      </div>

    </>
  );
}
