package xyz.wmmp.gallery.server.controllers;

import java.io.IOException;

import org.springframework.data.domain.Slice;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import xyz.wmmp.gallery.server.data.ThumbnailDTO;
import xyz.wmmp.gallery.server.data.Image;
import xyz.wmmp.gallery.server.data.ImageUploadMetadata;
import xyz.wmmp.gallery.server.services.ImageService;

@RestController
@RequestMapping("/images")
public class ImageController{

  private final ImageService imageService;

  public ImageController(ImageService imageService){
    this.imageService = imageService;
  }

  @GetMapping("/thumbnails")
  public Slice<ThumbnailDTO> getThumbnails(
    @RequestParam(defaultValue = "0") int page,
    @RequestParam(defaultValue = "20") int size,
    @RequestParam(defaultValue = "asc") String sort
  ){

    Sort.Direction direction = Sort.Direction.fromString(sort);
    return imageService.getThumbnails(page, size, direction);
    
  }

  @GetMapping("/{id}")
  public Image getImage(@PathVariable Long id){
    return imageService.getImage(id);
  }

  @PostMapping("/images")
  public ResponseEntity<Void> createImage(
    @RequestPart("file") MultipartFile file,
    @RequestPart("metadata") ImageUploadMetadata metadata
  ) throws IOException {
    imageService.createImage(file, metadata);
    return ResponseEntity.status(HttpStatus.CREATED).build();
  }

  @PutMapping("/images/{id}")
  public ResponseEntity<Void> updateImage(
    @PathVariable Long id,
    @RequestBody ImageUploadMetadata metadata
  ) throws IOException {
    imageService.updateImage(id, metadata);
    return ResponseEntity.noContent().build();
  }

  @DeleteMapping("/images/{id}")
  public ResponseEntity<Void> deleteImage(@PathVariable Long id){
    imageService.deleteImage(id);
    return ResponseEntity.noContent().build();
  }
  
  
  @ExceptionHandler(IllegalArgumentException.class)
  public ResponseEntity<String> handleBadSort(IllegalArgumentException ex){
    return ResponseEntity.badRequest().body("Invalid sort argument, use 'asc' or 'desc'");
  }
}
