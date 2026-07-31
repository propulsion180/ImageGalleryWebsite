package xyz.wmmp.gallery.server.controllers;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import xyz.wmmp.gallery.server.data.Image;
import xyz.wmmp.gallery.server.repositories.ImageRepository;

@RestController
public class SitemapController {
  private final ImageRepository imageRepository;

  @Autowired
  public SitemapController(ImageRepository imageRepository){
    this.imageRepository = imageRepository;
  }

  @GetMapping(value = "/sitemap.xml", produces = MediaType.APPLICATION_XML_VALUE)
  public String getSiteMap(){
    StringBuilder xml = new StringBuilder();

    xml.append("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n");
    xml.append("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n");

    xml.append(urlEntry("https://gallery.wmmp.xyz/", null));
    xml.append(urlEntry("https://gallery.wmmp.xyz/contactform", null));

    List<Image> images = imageRepository.findAll();
    for(Image image : images){
      String loc = "https://gallery.wmmp.xyz/single/" + image.getId();
      xml.append(urlEntry(loc, image.getUploadedTime()));
    }

    xml.append("</urlset>");
 
    return xml.toString();
  }

  public String urlEntry(String loc, LocalDateTime lastMod){
    StringBuilder entry = new StringBuilder();
    entry.append("  <url>\n");
    entry.append("    <loc>").append(loc).append("</loc>\n");
    if(lastMod != null){
      entry.append("    <lastmod>").append(lastMod.toLocalDate()).append("</lastmod>\n");
    }
    entry.append("  </url>\n");
    return entry.toString();
  }
}
