package it.valeriovaudi.i18nmessage.messages

import com.jayway.jsonpath.JsonPath
import org.springframework.jms.annotation.JmsListener
import org.springframework.messaging.rsocket.RSocketRequester
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono
import java.util.*

@Component
class AwsSQSListener(private val messageRepository: MessageRepository,
                     private val requester: Mono<RSocketRequester>) {

    @JmsListener(destination = "i18n-messages-updates")
    fun onMessage(message: String) {
        val s3key = JsonPath.read(message, "$.detail.requestParameters") as String
        val applicationName = s3key.split("//").first()
        messageRepository.find(applicationName, Locale.ENGLISH)
                .flatMap { bundle ->
                    requester.flatMap { req ->
                        req.route("messages.$applicationName")
                                .data(bundle)
                                .send()
                    }
                }.subscribe()

        println(message)
    }
}