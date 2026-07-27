package xyz.wmmp.gallery.server.data;

import java.time.Instant;
import java.util.List;

public record MonitoringDTO(
  long diskUsedBytes,
  long diskTotalBytes,
  long jvmUsedMemoryBytes,
  long jvmMaxMemoryBytes,
  long imageCount,
  long totalImageStorageBytes,
  List<String> recentErrors,
  List<LoginAttemptDTO> recentFailedLogins,
  Instant lastSuccessfulLogin,
  long requestCountLastHour,
  long errorCountLastHour
){}


