package xyz.wmmp.gallery.server.repositories;

import org.springframework.data.jpa.repository.JpaRepository;

import xyz.wmmp.gallery.server.data.User;

public interface UserRepository extends JpaRepository<User, Long>{

  public User findByUName(String uName);
    
}
