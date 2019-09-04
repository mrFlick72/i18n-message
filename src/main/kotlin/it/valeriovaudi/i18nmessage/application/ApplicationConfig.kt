package it.valeriovaudi.i18nmessage.application

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.cassandra.core.CassandraTemplate

@Configuration
class ApplicationConfig {

    @Bean
    fun applicationRepository(cassandraTemplate: CassandraTemplate) =
            CassandraApplicationRepository(cassandraTemplate)
}