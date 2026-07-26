package xyz.wmmp.gallery.server.services;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Component;

import xyz.wmmp.gallery.server.data.User;
import xyz.wmmp.gallery.server.data.UserType;
import xyz.wmmp.gallery.server.repositories.UserRepository;

@Component
public class OwnerSeeder implements CommandLineRunner{
  private final UserRepository userRepository;
  private final PasswordEncoder passwordEncoder;

  @Autowired
  public OwnerSeeder(UserRepository userRepository, PasswordEncoder passwordEncoder){
    this.userRepository = userRepository;
    this.passwordEncoder = passwordEncoder;
  }

  @Override
  public void run(String... args){
    if(userRepository.count() == 0){
      User owner = new User();
      owner.setUsername("testuser");
      owner.setPasswordHash(passwordEncoder.encode("testertester2020"));
      owner.setPerms(UserType.ADMIN);
      userRepository.save(owner);
      System.out.println("Using Seeded owner account test");
    }
  }
}
