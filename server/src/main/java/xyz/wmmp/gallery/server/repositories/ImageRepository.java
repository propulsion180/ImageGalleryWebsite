package xyz.wmmp.gallery.server.repositories;

import java.time.Instant;
import java.util.List;

import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Slice;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import xyz.wmmp.gallery.server.data.Image;
import xyz.wmmp.gallery.server.data.ThumbnailDTO;

public interface ImageRepository extends JpaRepository<Image, Long>{
  @Query("SELECT new xyz.wmmp.gallery.server.data.ThumbnailDTO(i.id, i.filename, i.uploadedTime) FROM Image i")
  Slice<ThumbnailDTO> findThumbnails(Pageable pageable);

  @Query("SELECT COALESCE(SUM(i.fileSizeBytes), 0) FROM Image i")
  long sumFileSizeBytes();

  long countByUploadedTimeAfter(Instant time);
}
