package xyz.wmmp.gallery.server.controllers;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.JwtException;
import jakarta.servlet.http.HttpServletResponse;

import org.apache.catalina.connector.Response;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.CookieValue;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.client.ResponseExtractor;
import org.springframework.web.server.ResponseStatusException;

import xyz.wmmp.gallery.server.authsec.CustomUserDetailsService;
import xyz.wmmp.gallery.server.authsec.JwtUtil;
import xyz.wmmp.gallery.server.data.User;
import xyz.wmmp.gallery.server.data.UserDTO;
import xyz.wmmp.gallery.server.repositories.UserRepository;

import java.time.Duration;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.UUID;

@Controller
public class AuthController {

    @Autowired private CustomUserDetailsService userDetailsService;
    @Autowired private PasswordEncoder passwordEncoder;
    @Autowired private JwtUtil jwtUtil;
    @Autowired private UserRepository userRepository;

    @PostMapping("/login")
    public ResponseEntity<UserDTO> login(@RequestBody LoginRequest request, HttpServletResponse response){
        UserDetails userDetails = userDetailsService.loadUserByUsername(request.uName());

        if(!passwordEncoder.matches(request.plainPassword, userDetails.getPassword())){ //check password
            throw new BadCredentialsException("Invalid credentials");
        }

        User user = userRepository.findByUName(request.uName());

        String jti = UUID.randomUUID().toString();
        String token = jwtUtil.generateToken(jti, user.getId().toString(), user.getPerms().toString());

        user.setJtiToken(jti);
        user.setTokenExpiry(Instant.now().plus(1, ChronoUnit.DAYS));
        userRepository.save(user);

        ResponseCookie cookie = ResponseCookie.from("session", token)
                .httpOnly(true)
                .secure(true)
                .sameSite("Strict")
                .maxAge(Duration.ofDays(1))
                .path("/")
                .build();
        response.addHeader(HttpHeaders.SET_COOKIE, cookie.toString());

        return ResponseEntity.ok(UserDTO.from(user));
    }

    @PreAuthorize("isAuthenticated()")
    @PostMapping("/tknlgn")
    public ResponseEntity<UserDTO> tokenLogin(@CookieValue(value = "session", required = true) String token ,HttpServletResponse response){
        if (token != null){
            try{
                Claims claims = jwtUtil.validate(token);
                String userId = claims.getSubject();
                User user = userRepository.findById(Long.getLong(userId)).orElseThrow();
                return ResponseEntity.ok(UserDTO.from(user));
            }catch (JwtException ignored){}
        }

        throw new ResponseStatusException(HttpStatus.NOT_FOUND, "You are not logged in yet! :)");
    }

    @PreAuthorize("isAuthenticated()")
    @PostMapping("/logout")
    public ResponseEntity<Void> logout(@CookieValue(value = "session", required = false) String token, HttpServletResponse response){
        if (token != null){
            try{
                Claims claims = jwtUtil.validate(token);
                String userId = claims.getSubject();
                User user = userRepository.findById(Long.getLong(userId)).orElseThrow();
                user.setJtiToken(null);
                user.setTokenExpiry(null);
                userRepository.save(user);
            }catch (JwtException ignored){}
        }

        ResponseCookie clear = ResponseCookie.from("session", "")
                .maxAge(0).path("/").build();
        response.addHeader(HttpHeaders.SET_COOKIE, clear.toString());
        return ResponseEntity.accepted().build();
    }

    @PreAuthorize(("hasRole('ADMIN')"))
    @PostMapping("/admin/change-password")
    public ResponseEntity<Void> changePassword(@CookieValue(value = "session", required = false) String token, @RequestBody ChangePasswordRequest request, HttpServletResponse response){
        if(request.newPassword() == null || request.newPassword().isBlank()){ return ResponseEntity.badRequest().build(); }
        if (token != null){
            try{
                Claims claims = jwtUtil.validate(token);
                String userId = claims.getSubject();
                User user = userRepository.findById(Long.getLong(userId)).orElseThrow();
                user.setJtiToken(null);
                user.setTokenExpiry(null);
                user.setPasswordHash(passwordEncoder.encode(request.newPassword()));
            }catch (JwtException ignored){}
        }

        ResponseCookie clear = ResponseCookie.from("session", "")
                .maxAge(0).path("/").build();
        response.addHeader(HttpHeaders.SET_COOKIE, clear.toString());

        return ResponseEntity.accepted().build();
    }

    record LoginRequest(String uName, String plainPassword){}
    record ChangePasswordRequest(String newPassword){}
}
