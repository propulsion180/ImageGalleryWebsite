package xyz.wmmp.gallery.server.trackers;

import java.time.Instant;
import java.util.Deque;
import java.util.List;
import java.util.concurrent.ConcurrentLinkedDeque;

import org.springframework.stereotype.Component;

@Component
public class RecentErrorTracker {
  private final Deque <String> recentErrors = new ConcurrentLinkedDeque<>();
  private static final int MAX_ENTRIES = 100;

  public void record(String message){
    recentErrors.addFirst(Instant.now() + " - " + message);
    while(recentErrors.size() > MAX_ENTRIES){
      recentErrors.removeLast();
    }
  }

  public List<String> getRecent(int limit){
    return recentErrors.stream().limit(limit).toList();
  }
}
