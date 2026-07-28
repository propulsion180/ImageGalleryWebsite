package xyz.wmmp.gallery.server.authsec;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.FormLoginConfigurer;
import org.springframework.security.config.annotation.web.configurers.LogoutConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

import xyz.wmmp.gallery.server.trackers.MetricFilter;

@Configuration
@EnableWebSecurity
public class SecurityConfig {
    @Autowired private JwtAuthFilter jwtAuthFilter;
    @Autowired private MetricFilter metricFilter;

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
                .csrf(csrf -> csrf.disable())
                .sessionManagement(s -> s
                        .sessionCreationPolicy(SessionCreationPolicy.STATELESS)
                )
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/", "/index.html", "/*.js", "/*.css", "/favicon.ico", "/login").permitAll()
                        .requestMatchers("/images/files/**").permitAll()
                        .requestMatchers(HttpMethod.GET, "/images/**").permitAll()
                        .requestMatchers(HttpMethod.POST, "/images/**").hasRole("ADMIN")
                        .requestMatchers(HttpMethod.PUT, "/images/**").hasRole("ADMIN")
                        .requestMatchers("/images/all").hasRole("ADMIN")
                        .requestMatchers(HttpMethod.POST, "/contact/**").permitAll()
                        .requestMatchers("/contact/**").hasRole("ADMIN")
                        
                        .requestMatchers("/admin/metrics/**").hasRole("ADMIN")
                        .requestMatchers("/admin/**").permitAll()
                        
                        .anyRequest().authenticated()
                )
                .addFilterBefore(jwtAuthFilter, UsernamePasswordAuthenticationFilter.class)
                .addFilterBefore(metricFilter, JwtAuthFilter.class);
        return http.build();
    }
}
