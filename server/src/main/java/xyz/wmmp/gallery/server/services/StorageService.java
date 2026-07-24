package xyz.wmmp.gallery.server.services;

import java.io.IOException;

import org.springframework.core.io.Resource;

public interface StorageService {
  void save (String relativePath, byte[] data) throws IOException;
  Resource load(String relativePath) throws IOException;
  void delete(String relativePath) throws IOException;
}
