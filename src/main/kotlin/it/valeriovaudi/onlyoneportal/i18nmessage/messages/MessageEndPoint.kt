package it.valeriovaudi.onlyoneportal.i18nmessage.messages

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.reactive.function.BodyInserters
import org.springframework.web.reactive.function.server.router
import java.util.*

@Configuration
class MessageEndPoint(private val messageRepository: MessageRepository) {

    @Bean
    fun messageEndPointRoute() =
            router {
                GET("/messages/{application}") {
                    messageRepository.find(
                            it.pathVariable("application"),
                            it.queryParam("lang").map { Locale(it) }.orElse(Locale.ENGLISH)
                    ).flatMap { messages -> ok().body(BodyInserters.fromValue(messages)) }
                }
            }
}