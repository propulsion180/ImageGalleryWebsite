package xyz.wmmp.gallery.server.data;

import java.net.FileNameMap;
import java.time.LocalDateTime;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.FetchType;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
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
  
  private String contentType;

  private Long fileSizeBytes;

  private String camera;

  @Enumerated(EnumType.STRING)
  @Column(nullable = false)
  private ImageType type;

  private String aperture;

  private String shutterSpeed;

  @Column(nullable = false)
  private Integer iso;

  private String filmStock;

  private String location;

  private String description;

  private LocalDateTime uploadedTime;

}
