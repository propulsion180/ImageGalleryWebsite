package xyz.wmmp.gallery.server.controllers;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import xyz.wmmp.gallery.server.data.ContactRequest;
import xyz.wmmp.gallery.server.repositories.ContactRequestRepository;

@RestController
@RequestMapping("/contact")
public class ContactController{

  private final ContactRequestRepository contactRequestRepository;

  @Autowired
  public ContactController(ContactRequestRepository contactRequestRepository){
    this.contactRequestRepository = contactRequestRepository;
  }

  @PostMapping
  public ResponseEntity<Void> submitContact(@RequestBody ContactRequest request){
    if(request.getId() != -1){return ResponseEntity.status(HttpStatus.BAD_REQUEST).build();}
    ContactRequest cr = new ContactRequest();
    cr.setName(request.getName());
    cr.setEmail(request.getEmail());
    cr.setMessage(request.getMessage());
    cr.setReferenceImageId(request.getReferenceImageId());
    contactRequestRepository.save(cr);
    return ResponseEntity.accepted().build();
  }

  @PreAuthorize("hasRole('ADMIN')")
  @GetMapping
  public List<ContactRequest> getAllContactRequests(){
    return contactRequestRepository.findAll();
  }

  @PreAuthorize("hasRole('ADMIN')")
  @DeleteMapping("/{id}")
  public ResponseEntity<Void> deleteContactReques(@PathVariable Long id){
    this.contactRequestRepository.deleteById(id);
    return ResponseEntity.noContent().build();    
  }
  
}
