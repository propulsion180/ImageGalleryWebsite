package xyz.wmmp.gallery.server.trackers;

import java.io.IOException;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@Component
public class MetricFilter extends OncePerRequestFilter{
  private final RequestMetricsTracker metricsTracker;

  @Autowired
  public MetricFilter(RequestMetricsTracker requestMetricsTracker){
    this.metricsTracker = requestMetricsTracker;
  }

  @Override
  protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain chain) throws ServletException, IOException{
    chain.doFilter(request, response);
    boolean isError = response.getStatus() >= 400;
    metricsTracker.recordRequest(isError);
  } 
}
