package it.valeriovaudi.onlyoneportal.i18nmessage.messages

import com.jayway.jsonpath.JsonPath
import org.springframework.jms.annotation.JmsListener
import org.springframework.messaging.rsocket.RSocketRequester
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono
import java.util.*

@Component
class AwsSQSListener(private val messageRepository: MessageRepository,
                     private val requesters: Map<String, Mono<RSocketRequester>>) {

    @JmsListener(destination = "i18n-messages-updates")
    fun onMessage(message: String) {
        applicationNameFor(message)
                .map { applicationName ->

                    println("application $applicationName bundle are refreshing")

                    messageRepository.find(applicationName, Locale.ENGLISH)
                            .flatMap { bundle ->
                                sendBundleToClient(bundle, applicationName)
                            }.subscribe()
                }
    }

    private fun applicationNameFor(message: String) =
            Optional.ofNullable(JsonPath.read(message, "$.application.value") as String?)
                    .map { key -> key.split("/").first() }


    private fun sendBundleToClient(bundle: Messages, applicationName: String): Mono<Void>? {
        return Optional.ofNullable(bundle)
                .map { requesters[applicationName] }
                .map { applicationRequester ->
                    applicationRequester?.flatMap { req ->
                        req.route("messages.$applicationName")
                                .data(bundle)
                                .send()
                    }
                }.orElse(Mono.empty())
    }
}