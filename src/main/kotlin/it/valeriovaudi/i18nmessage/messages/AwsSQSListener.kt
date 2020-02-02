package it.valeriovaudi.i18nmessage.messages

import com.jayway.jsonpath.JsonPath
import org.springframework.jms.annotation.JmsListener
import org.springframework.messaging.rsocket.RSocketRequester
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono
import java.util.*
import java.util.regex.Pattern
import kotlin.collections.LinkedHashMap

@Component
class AwsSQSListener(private val messageRepository: MessageRepository,
                     private val requester: Mono<RSocketRequester>) {

    @JmsListener(destination = "i18n-messages-updates")
    fun onMessage(message: String) {
        Optional.ofNullable(JsonPath.read(message, "$.detail.requestParameters.key") as String?)
                .map {
                    val applicationName = it.split("/").first()
                    println("application $applicationName bundle are refreshing")
                    messageRepository.find(applicationName, Locale.ENGLISH)
                            .flatMap { bundle ->
                                Optional.ofNullable(bundle)
                                        .map {
                                            requester.flatMap { req ->
                                                req.route("messages.$applicationName")
                                                        .data(bundle)
                                                        .send()
                                            }
                                        }.orElse(Mono.empty())
                            }.subscribe()
                }
    }
}