package xyz.wmmp.gallery.server.services;

import java.awt.Graphics2D;
import java.awt.RenderingHints;
import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.time.LocalDateTime;
import java.util.List;
import java.util.NoSuchElementException;

import javax.imageio.IIOImage;
import javax.imageio.ImageIO;
import javax.imageio.ImageWriteParam;
import javax.imageio.ImageWriter;
import javax.imageio.spi.ImageReaderWriterSpi;
import javax.imageio.stream.ImageOutputStream;
import javax.print.attribute.standard.MediaPrintableArea;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Slice;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.servlet.mvc.HttpRequestHandlerAdapter;

import com.drew.imaging.ImageMetadataReader;
import com.drew.metadata.Metadata;
import com.drew.metadata.exif.ExifIFD0Directory;

import ch.qos.logback.classic.joran.sanity.IfNestedWithinSecondPhaseElementSC;
import xyz.wmmp.gallery.server.data.Image;
import xyz.wmmp.gallery.server.data.ImageUploadMetadata;
import xyz.wmmp.gallery.server.data.ThumbnailDTO;
import xyz.wmmp.gallery.server.repositories.ImageRepository;
import xyz.wmmp.gallery.server.data.ImageType;
import xyz.wmmp.gallery.server.services.StorageService;

@Service
public class ImageService {

  private final ImageRepository imageRepository;
  private final StorageService storageService;

  @Autowired
  public ImageService(ImageRepository imageRepository, StorageService storageService){
    this.imageRepository = imageRepository;
    this.storageService = storageService;
  }

  public Slice<ThumbnailDTO> getThumbnails(int page, int size, Sort.Direction direction){
    Pageable pageable = PageRequest.of(page, size, Sort.by(direction, "uploadedTime"));
     return imageRepository.findThumbnails(pageable);
  }

  public Image getImage(Long id){
    return imageRepository.findById(id).orElseThrow(() -> new NoSuchElementException("An image with the id " + id + " was not found in the repository"));
  }

  public List<Image> all(){
    return imageRepository.findAll();
  }

  public void createImage(MultipartFile file, ImageUploadMetadata metadata) throws IOException {
    System.out.println("Entered createImage");
    BufferedImage original;
    try(InputStream is = file.getInputStream()){
      original = ImageIO.read(is);
    }

    if(original == null){throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Unsupported or corrupt image file");}

    int orientation = getExifOrientation(file.getBytes());
    original = applyOrientation(original, orientation);

    String fileName = file.getOriginalFilename();
    byte[] fullResBytes = compressToJpg(original, 1.0f);
    storageService.save("full/" + fileName, fullResBytes);
    System.out.println("Read file its name is " + fileName);
    System.out.println("Saving original");
    BufferedImage thumbnail = resize(original, 400);
    byte[] thumbBytes = compressToJpg(thumbnail, 0.6f);
    storageService.save("thumb/" + fileName, thumbBytes);
    System.out.println("Savign thumb");
    Image i = new Image();
    i.setFilename(fileName);
    i.setContentType("image/jpeg");
    i.setFileSizeBytes((long) fullResBytes.length);
    i.setCamera(metadata.camera());
    if(metadata.filmStock().isEmpty() || metadata.filmStock() == null){
      i.setAperture(metadata.aperture());
      i.setShutterSpeed(metadata.shutterSpeed());
      i.setType(ImageType.DIGITAL);
    }else{
      i.setType(ImageType.FILM);
      i.setFilmStock(metadata.filmStock());
    }
    i.setIso(metadata.iso());
    i.setLocation(metadata.location());
    i.setDescription(metadata.description());
    i.setUploadedTime(LocalDateTime.now());

    imageRepository.save(i);
    System.out.println("Saving image object");
  }

  public void updateImage(Long id, ImageUploadMetadata metadata) throws IOException {
    Image i = imageRepository.findById(id).orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "No image with id " + id));
    i.setCamera(metadata.camera());
    if(metadata.filmStock().isEmpty() || metadata.filmStock() == null){
      i.setAperture(metadata.aperture());
      i.setShutterSpeed(metadata.shutterSpeed());
      i.setType(ImageType.DIGITAL);
    }else{
      i.setType(ImageType.FILM);
      i.setFilmStock(metadata.filmStock());
    }
    i.setIso(metadata.iso());
    i.setLocation(metadata.location());
    i.setDescription(metadata.description());
    imageRepository.save(i); 
  }

  public void deleteImage(Long id){
    imageRepository.deleteById(id);
  }

  private byte[] compressToJpg(BufferedImage image, float quality) throws IOException{
    ByteArrayOutputStream baos = new ByteArrayOutputStream();
    ImageWriter iw = ImageIO.getImageWritersByFormatName("jpg").next();
    ImageWriteParam param = iw.getDefaultWriteParam();

    try(ImageOutputStream ios = ImageIO.createImageOutputStream(baos)){
      iw.setOutput(ios);
      iw.write(null, new IIOImage(image, null, null), param);
    }finally{
      iw.dispose();
    }
    return baos.toByteArray();
  }

  private BufferedImage resize(BufferedImage original, int maxWidth){
    int newWidth = Math.min(original.getWidth(), maxWidth);
    int newHeight = (int)((double) original.getHeight() / original.getWidth() * newWidth);
    BufferedImage resized = new BufferedImage(newWidth, newHeight, BufferedImage.TYPE_INT_RGB);
    Graphics2D g = resized.createGraphics();
    g.setRenderingHint(RenderingHints.KEY_INTERPOLATION, RenderingHints.VALUE_INTERPOLATION_BILINEAR);
    g.drawImage(original, 0, 0, newWidth, newHeight, null);
    g.dispose();
    return resized;
  }

  private int getExifOrientation(byte[] imageBytes){
    try{
      Metadata metadata = ImageMetadataReader.readMetadata(new ByteArrayInputStream(imageBytes));
      ExifIFD0Directory directory = metadata.getFirstDirectoryOfType(ExifIFD0Directory.class);
      if(directory != null && directory.containsTag(ExifIFD0Directory.TAG_ORIENTATION)){
        return directory.getInt(ExifIFD0Directory.TAG_ORIENTATION);
      }
    }catch(Exception e){
      
    }
    return 1;
  }

  private BufferedImage applyOrientation(BufferedImage image, int orientation){
    return switch(orientation){
      case 3 -> rotate(image, 180);
      case 6 -> rotate(image, 90);
      case 8 -> rotate(image, 270);
      default -> image;
    };
  }

  private BufferedImage rotate(BufferedImage image, int degrees){
    double radians = Math.toRadians(degrees);
    int w = image.getWidth(), h = image.getHeight();

    int newW = (degrees == 90 || degrees == 270)? h : w;
    int newH = (degrees == 90 || degrees == 270)? w : h;

    BufferedImage rotated = new BufferedImage(newW, newH, image.getType());
    Graphics2D g = rotated.createGraphics();
    g.setRenderingHint(RenderingHints.KEY_INTERPOLATION, RenderingHints.VALUE_INTERPOLATION_BILINEAR);
    g.translate((newW - w) / 2, (newH - h) / 2);
    g.rotate(radians, w / 2.0, h / 2.0);
    g.drawImage(image, 0, 0, null);
    g.dispose();
    return rotated;
  }
}
