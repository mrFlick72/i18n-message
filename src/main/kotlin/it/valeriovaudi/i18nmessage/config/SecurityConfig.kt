package it.valeriovaudi.i18nmessage.config

import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter

@EnableWebSecurity
class SecurityConfig : WebSecurityConfigurerAdapter() {

    @Throws(Exception::class)
    override fun configure(http: HttpSecurity) {
        http.csrf().disable()
                .authorizeRequests().mvcMatchers("/actuator/**").permitAll().and()
                .authorizeRequests().anyRequest().authenticated().and()
                .oauth2ResourceServer().jwt()
    }

}
