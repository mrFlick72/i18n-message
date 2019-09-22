package it.valeriovaudi.i18nmessage.messages

import org.springframework.http.ResponseEntity.ok
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController
import java.util.*

@RestController
class MessageEndPoint(private val messageRepository: MessageRepository) {

    @GetMapping("/messages/{application}")
    fun findAllMessages(@PathVariable("application") application: String,
                        @RequestParam("lang", required = false, defaultValue = "en") locale: Locale) =
    ok(messageRepository.find(application, locale))
}