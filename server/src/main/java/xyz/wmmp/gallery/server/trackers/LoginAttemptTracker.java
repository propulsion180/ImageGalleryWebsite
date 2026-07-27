package xyz.wmmp.gallery.server.trackers;

import java.time.Instant;
import java.util.Deque;
import java.util.List;
import java.util.concurrent.ConcurrentLinkedDeque;

import org.springframework.stereotype.Component;

import xyz.wmmp.gallery.server.data.LoginAttemptDTO;

@Component
public class LoginAttemptTracker{
  private final Deque<LoginAttemptDTO> failedAttempts = new ConcurrentLinkedDeque<>();
  private volatile Instant lastSuccessfulLogin;
  private static final int MAX_ENTRIES = 50;

  public void recordFailure(String username){
    failedAttempts.addFirst(new LoginAttemptDTO(username, Instant.now()));
    while(failedAttempts.size() > MAX_ENTRIES){failedAttempts.removeLast();}
  }

  public void recordSuccess(){
    lastSuccessfulLogin = Instant.now();
  }

  public List<LoginAttemptDTO> getRecentFailures(int limit){
    return failedAttempts.stream().limit(limit).toList();
  }

  public Instant getLastSuccessfulLogin(){
    return lastSuccessfulLogin;
  }
}

