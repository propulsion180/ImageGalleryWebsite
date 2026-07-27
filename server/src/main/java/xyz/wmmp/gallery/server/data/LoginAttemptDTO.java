package xyz.wmmp.gallery.server.data;

import java.time.Instant;

public record LoginAttemptDTO(String username, Instant timestamp){}
