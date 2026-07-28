package xyz.wmmp.gallery.server.repositories;

import org.springframework.data.jpa.repository.JpaRepository;
import xyz.wmmp.gallery.server.data.ContactRequest;

public interface ContactRequestRepository extends JpaRepository<ContactRequest, Long>{
  
}
