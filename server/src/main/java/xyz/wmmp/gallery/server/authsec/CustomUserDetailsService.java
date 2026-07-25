package xyz.wmmp.gallery.server.authsec;

import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Component;

import xyz.wmmp.gallery.server.repositories.UserRepository;

@Component
public class CustomUserDetailsService implements UserDetailsService {
    private final UserRepository userRepository;

    public CustomUserDetailsService(UserRepository userRepository){
        this.userRepository = userRepository;
    }

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        xyz.wmmp.gallery.server.data.User user = userRepository.findByUName(username);

        return User.builder()
                .username(user.getUName())
                .password(user.getPasswordHash())
                .roles(user.getPerms().toString())
                .build();
    }
}
