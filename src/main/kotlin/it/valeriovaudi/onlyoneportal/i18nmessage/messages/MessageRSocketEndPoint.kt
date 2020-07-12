package it.valeriovaudi.onlyoneportal.i18nmessage.messages

import org.springframework.messaging.handler.annotation.DestinationVariable
import org.springframework.messaging.handler.annotation.MessageMapping
import org.springframework.stereotype.Controller
import java.util.*

@Controller
class MessageRSocketEndPoint(private val messageRepository: MessageRepository) {

    @MessageMapping("messages.{application}")
    fun messageEndPointRoute(
            @DestinationVariable("application") application: String,
            lang: String?
    ) = messageRepository.find(application, Optional.ofNullable(lang).map { Locale(it) }.orElse(Locale.ENGLISH))
}