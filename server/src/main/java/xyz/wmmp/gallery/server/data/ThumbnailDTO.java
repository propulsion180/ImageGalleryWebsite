package xyz.wmmp.gallery.server.data;

import java.time.LocalDateTime;

public record ThumbnailDTO(Long id, String path, LocalDateTime uploadedTime){
  public static ThumbnailDTO from(Long id, String fileName, LocalDateTime uploadedTime){
    return new ThumbnailDTO(id, "thumb/" + fileName, uploadedTime);
  }
}
