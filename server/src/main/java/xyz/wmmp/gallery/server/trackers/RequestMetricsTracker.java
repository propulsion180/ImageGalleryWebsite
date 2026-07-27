package xyz.wmmp.gallery.server.trackers;

import java.time.Duration;
import java.time.Instant;
import java.util.concurrent.atomic.AtomicLong;

import org.springframework.stereotype.Component;

@Component
public class RequestMetricsTracker{
  private final AtomicLong requestCount = new AtomicLong();
  private final AtomicLong errorCount = new AtomicLong();
  private volatile Instant windowStart = Instant.now();
  
  public void recordRequest(boolean isError){
    resetifWindowExpired();
    requestCount.incrementAndGet();
    if(isError){errorCount.incrementAndGet();}
  }

  private void resetifWindowExpired(){
    if(Duration.between(windowStart, Instant.now()).toHours() >= 1){
      requestCount.set(0);
      errorCount.set(0);
      windowStart = Instant.now();
    }
  }

  public long getRequestCount(){ return requestCount.get();}
  public long getErrorCount(){ return errorCount.get();}
  
}
