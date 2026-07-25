package xyz.wmmp.gallery.server.data;

public record UserDTO(Long id, String userName, UserType perms){
  public static UserDTO from(User u){
    return new UserDTO(u.getId(), u.getUName(), u.getPerms());
  }
}
