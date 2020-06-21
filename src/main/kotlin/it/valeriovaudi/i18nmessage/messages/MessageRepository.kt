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
            client.baseUrl(baseUrlFor(application))
                    .build()
                    .get()
                    .accept(MediaType.APPLICATION_OCTET_STREAM)
                    .exchange()
                    .flatMap { it.bodyToMono(ByteArray::class.java) }
                    .flatMap { loadBundle(it) }


    private fun loadBundle(data: ByteArray): Mono<Map<String, String>> =
            Mono.defer { Properties().toMono() }
                    .map { props ->
                        ByteArrayInputStream(data).use { props.load(it) };
                        props
                    }
                    .map { it.toMap() as Map<String, String> }

    private fun baseUrlFor(application: String) =
            "$repositoryServiceUrl/documents/$registrationName?path=$application&fileName=messages&fileExt=properties"
}

typealias  S3ObjectSummaryPredicate = (S3ObjectSummary) -> Boolean

open class AwsS3MessageRepository(private val s3client: AmazonS3,
                                  private val bucketName: String) : MessageRepository {

    override fun find(application: String, language: Locale): Mono<Messages> =
            getAllS3MessagesBundleFor(application)
                    .flatMap { s3MessageBundleFor(it, application, language.toString()) }
                    .flatMap { s3MessageBundleContentFor(it) }

    private fun s3MessageBundleFor(message: ObjectListing, application: String, language: String) =
            Mono.defer {
                Optional.ofNullable(message.objectSummaries.find(resourceBundleFinderPredicate(application, language)))
                        .orElse(message.objectSummaries.find(defaultResourceBundleFinderPredicate(application)))
                        .toMono()
            }

    private fun defaultResourceBundleFinderPredicate(messagesKey: String): S3ObjectSummaryPredicate =
            { it.key == "$messagesKey/messages.properties" }

    private fun resourceBundleFinderPredicate(messagesKey: String, language: String): S3ObjectSummaryPredicate =
            { it.key == "$messagesKey/messages_$language.properties" }

    private fun s3MessageBundleContentFor(message: S3ObjectSummary) =
            Mono.fromCallable {
                s3client.getObject(bucketName, message.key)
                        .objectContent.buffered()
            }.flatMap { loadBundle(it) }

    private fun loadBundle(data: BufferedInputStream): Mono<Map<String, String>> =
            Mono.defer { Properties().toMono() }
                    .map { it.load(data); it }
                    .map { it.toMap() as Map<String, String> }


    private fun getAllS3MessagesBundleFor(application: String) = Mono.fromCallable { s3client.listObjects(bucketName, "$application") }

}
