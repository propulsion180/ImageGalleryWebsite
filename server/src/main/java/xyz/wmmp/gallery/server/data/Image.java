package xyz.wmmp.gallery.server.data;

import java.net.FileNameMap;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
@Entity
@Table(name = "Images")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Image {

  @Id
  @GeneratedValue(strategy = GenerationType.IDENTITY)
  private Long id;

  private String filename;

  private String thumbnailFilename;

  private String contentType;

  private Long fileSizeBytes;

  private String camera

  private 

}
