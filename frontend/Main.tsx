import React, { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ImageData, User } from "./App";

interface MainProps {
  user: User | null;
}

interface ThumbnailDTO {
  id: number;
  path: string;
  uploadDateTime: string;
}

interface SliceResponse<T> {
  content: T[];
  last: boolean;
}

type SortDirection = "asc" | "desc";

async function fetchThumbnails(page: number, sort: SortDirection): Promise<SliceResponse<ThumbnailDTO>> {
  const response = await fetch(`/images/thumbnails?page=${page}&size=20&sort=${sort}`, {
    credentials: "include",
  });

  if(!response.ok) throw new Error('Failed to fetch thumbnails: ${res.status}');
  return response.json();
}


function InfiniteScrollTrigger({ onIntersect }: { onIntersect: () => void}){
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if(entries[0].isIntersecting){
          onIntersect();
        }
      },
      { threshold: 1.0}
    );

    const el = ref.current;
    if(el) observer.observe(el);

    return () => {
      if(el) observer.unobserve(el);
    };
  }, [onIntersect]);

  return <div ref={ref} style={{ height: "1px" }} />;
}


export default function Main({ user }: MainProps) {
  const navigate = useNavigate();

  const [thumbnails, setThumbnails] = useState<ThumbnailDTO[]>([]);
  const [page, setPage] = useState(0);
  const [sort, setSort] = useState<SortDirection>("asc");
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  async function loadPage(pageToLoad: number, currentSort: SortDirection, isReset: boolean){
    setLoading(true);
    try{
      const data = await fetchThumbnails(pageToLoad, currentSort);
      setThumbnails(prev => isReset ? data.content: [...prev, ...data.content]);
      setHasMore(!data.last);
      setPage(pageToLoad);
    }finally{
      setLoading(false);
    }
  }

  
  useEffect(() => {
    setThumbnails([]);
    setPage(0);
    setHasMore(true);
    loadPage(0, sort, true);
  }, [sort]);

  function loadMore(){
    if(!loading && hasMore){
      loadPage(page + 1, sort, false);
    }
  }

  function toggleSort(){
    setSort(prev => (prev === "desc"? "asc": "desc"));
  }  

  return (
    <div>
      <button className="small-button" onClick={toggleSort}>
        Sort: {sort === "asc" ? "Chronological": "Reverse Chronological"}
      </button>
      <div className="allimage-grid">

        {thumbnails.map(thumbnail =>(
          
          <div className="allimage-item">
            <Link to={`/single/${thumbnail.id}`} >
              <img key= {thumbnail.id} src={"/images/files/thumb/" + thumbnail.path} alt={thumbnail.uploadDateTime} />
            </Link>
          </div>
        ))}
        {hasMore && (<InfiniteScrollTrigger onIntersect={loadMore} />)}
      </div>
    </div>
  );
}
