package xyz.wmmp.gallery.server.repositories;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import xyz.wmmp.gallery.server.data.User;

public interface UserRepository extends JpaRepository<User, Long>{

  public Optional<User> findByUsername(String username);
    
}
