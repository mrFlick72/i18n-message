package it.valeriovaudi.i18nmessage.messages

import com.amazonaws.services.s3.AmazonS3
import com.amazonaws.services.s3.model.ObjectListing
import com.amazonaws.services.s3.model.S3ObjectSummary
import reactor.core.publisher.Mono
import reactor.core.publisher.toMono
import java.io.BufferedInputStream
import java.util.*

typealias  S3ObjectSummaryPredicate = (S3ObjectSummary) -> Boolean

open class AwsS3MessageRepository(private val s3client: AmazonS3,
                                  private val bucketName: String) : MessageRepository {

    //    @Cacheable("i18nMessageBundle", key = "#application + '_' + #language.toString()")
    override fun find(application: String, language: Locale): Mono<Messages> =
            getAllS3MessagesBundleFor(application)
                    .flatMap { s3MessageBundleFor(it, application, language.toString()) }
                    .flatMap { s3MessageBundleContentFor(it) }

    private fun s3MessageBundleFor(message: ObjectListing, application: String, language: String) =
            Mono.defer {
                messageKeyFor(application)
                        .let { messagesKey ->
                            Optional.ofNullable(message.objectSummaries.find(resourceBundleFinderPredicate(messagesKey, language)))
                                    .orElse(message.objectSummaries.find(defaultResourceBundleFinderPredicate(messagesKey)))
                        }.toMono()
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


    private fun messageKeyFor(application: String) = "i18n-messages/$application"

    private fun getAllS3MessagesBundleFor(application: String) = Mono.fromCallable { s3client.listObjects(bucketName, messageKeyFor(application)) }

}