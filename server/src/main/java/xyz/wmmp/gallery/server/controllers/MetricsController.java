package xyz.wmmp.gallery.server.controllers;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.Instant;
import java.time.temporal.ChronoUnit;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import xyz.wmmp.gallery.server.data.MonitoringDTO;
import xyz.wmmp.gallery.server.repositories.ImageRepository;
import xyz.wmmp.gallery.server.trackers.LoginAttemptTracker;
import xyz.wmmp.gallery.server.trackers.RecentErrorTracker;
import xyz.wmmp.gallery.server.trackers.RequestMetricsTracker;

@RestController
@RequestMapping("/admin/metrics")
public class MetricsController{
  private final ImageRepository imageRepository;
  private final RecentErrorTracker recentErrorTracker;
  private final LoginAttemptTracker loginAttemptTracker;
  private final RequestMetricsTracker requestMetricsTracker;
  private final Path storageBasePath;

  public MetricsController(ImageRepository imageRepository, RecentErrorTracker recentErrorTracker, LoginAttemptTracker loginAttemptTracker, RequestMetricsTracker requestMetricsTracker, @Value("${storage.base-path}") String storageBasePath){
    this.imageRepository = imageRepository;
    this.recentErrorTracker = recentErrorTracker;
    this.loginAttemptTracker = loginAttemptTracker;
    this .requestMetricsTracker = requestMetricsTracker;
    this.storageBasePath = Paths.get(storageBasePath);
  }

  @PreAuthorize("hasRole('ADMIN')")
  @GetMapping("/stats")
  public MonitoringDTO getStats() throws IOException{
    File storageDir = storageBasePath.toFile();
    long diskUsed = storageDir.getTotalSpace() - storageDir.getFreeSpace();
    long diskTotal = storageDir.getTotalSpace();

    Runtime rt = Runtime.getRuntime();

    return new MonitoringDTO(
      diskUsed,
      diskTotal,
      rt.totalMemory() - rt.freeMemory(),
      rt.maxMemory(),
      imageRepository.count(),
      sumImageStorage(),
      recentErrorTracker.getRecent(20),
      loginAttemptTracker.getRecentFailures(20),
      loginAttemptTracker.getLastSuccessfulLogin(),
      requestMetricsTracker.getRequestCount(),
      requestMetricsTracker.getErrorCount()      
    );
  }

  private long sumImageStorage(){
    return imageRepository.sumFileSizeBytes();
  }
}
