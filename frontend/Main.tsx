import React from 'react'
import {ImageData} from './App'

interface MainProps {
    images: Map<string, ImageData>;
}



export default function Main({images}: MainProps) {
  return (
    <div>
        {Array.from(images).map(([key, val]) => (
            <img src={val.FilePath} alt={val.Description} className='image' />
        ))}
    </div>
  );
}
