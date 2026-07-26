package xyz.wmmp.gallery.server.services;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

@Service
public class LocalStorageService implements StorageService{
  
  private final Path basePath;

  public LocalStorageService(@Value("${storage.base-path}") String basePath){
    this.basePath = Paths.get(basePath);
  }

  @Override
  public void save(String relativePath, byte[] data) throws IOException{
    Path fullPath = resolveOrThrow(relativePath);
    Files.createDirectories(fullPath.getParent());
    Files.write(fullPath, data);
  }
  
  @Override
  public Resource load(String relativePath) throws IOException{
    Path fullPath = resolveOrThrow(relativePath);
    System.out.println("Loading somethign in localstorage: " + fullPath.toString());
    if(!Files.exists(fullPath)){
      throw new ResponseStatusException(HttpStatus.NOT_FOUND, "File not found: " + relativePath);
    }
    return new FileSystemResource(fullPath);
  }

  @Override
  public void delete(String relativePath) throws IOException {
    Files.deleteIfExists(resolveOrThrow(relativePath));
  }

  private Path resolveOrThrow(String relativePath){
    Path resolved = basePath.resolve(relativePath).normalize();
    if(!resolved.startsWith(basePath)){
      throw new IllegalArgumentException("Invalid path: " + relativePath);  
    }
    return resolved;
  }
  
}
