package it.valeriovaudi.i18nmessage.messages

import org.springframework.web.reactive.function.client.WebClient
import reactor.core.publisher.Mono
import reactor.core.publisher.toMono
import java.io.ByteArrayInputStream
import java.util.*

typealias Messages = Map<String, String>

interface MessageRepository {

    fun find(application: String, language: Locale): Mono<Messages>
}

class RestMessageRepository(
        private val repositoryServiceUrl: String,
        private val registrationName: String,
        private val client: WebClient.Builder) : MessageRepository {

    override fun find(application: String, language: Locale): Mono<Messages> =
            findFor { baseUrlFor(application, language.toString()) }
                    .onErrorResume { findFor { baseUrlFor(application) } }

    private fun findFor(baseURL: () -> String): Mono<Messages> =
            client.baseUrl(baseURL())
                    .build()
                    .get()
                    .exchange()
                    .flatMap {
                        when (it.statusCode().is2xxSuccessful) {
                            true -> it.bodyToMono(ByteArray::class.java)
                            false -> Mono.error(RuntimeException("language not found"))
                        }
                    }
                    .flatMap(this::loadBundle)


    private fun loadBundle(data: ByteArray): Mono<Map<String, String>> =
            Mono.defer { Properties().toMono() }
                    .map { props ->
                        ByteArrayInputStream(data).use { props.load(it) }
                        props
                    }
                    .map { it.toMap() as Map<String, String> }


    private fun baseUrlFor(application: String) =
            "$repositoryServiceUrl/documents/$registrationName?path=$application&fileName=messages&fileExt=properties"

    private fun baseUrlFor(application: String, lang: String) =
            "$repositoryServiceUrl/documents/$registrationName?path=$application&fileName=messages_$lang&fileExt=properties"
}