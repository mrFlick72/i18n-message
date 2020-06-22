package it.valeriovaudi.i18nmessage.messages

import com.amazonaws.services.s3.AmazonS3
import com.amazonaws.services.s3.model.ObjectListing
import com.amazonaws.services.s3.model.S3ObjectSummary
import org.springframework.http.MediaType
import org.springframework.web.reactive.function.client.WebClient
import reactor.core.publisher.Mono
import reactor.core.publisher.toMono
import java.io.BufferedInputStream
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
                    .switchIfEmpty(findFor { baseUrlFor(application) })

    private fun findFor(baseURL: () -> String): Mono<Messages> {
        return client.baseUrl(baseURL())
                .build()
                .get()
                .accept(MediaType.APPLICATION_OCTET_STREAM)
                .exchange()
                .flatMap { it.bodyToMono(ByteArray::class.java) }
                .flatMap { loadBundle(it) }
    }


    private fun loadBundle(data: ByteArray): Mono<Map<String, String>> =
            Mono.defer { Properties().toMono() }
                    .map { props ->
                        ByteArrayInputStream(data).use { props.load(it) };
                        props
                    }
                    .map { it.toMap() as Map<String, String> }

    private fun baseUrlFor(application: String) =
            "$repositoryServiceUrl/documents/$registrationName?path=$application&fileName=messages&fileExt=properties"

    private fun baseUrlFor(application: String, lang: String) =
            "$repositoryServiceUrl/documents/$registrationName?path=$application&fileName=messages_lang&fileExt=properties"
}