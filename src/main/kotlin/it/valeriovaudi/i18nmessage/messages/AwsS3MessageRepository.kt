package it.valeriovaudi.i18nmessage.messages

import com.amazonaws.services.s3.AmazonS3
import com.amazonaws.services.s3.model.ObjectListing
import com.amazonaws.services.s3.model.S3ObjectSummary
import java.io.BufferedInputStream
import java.util.*

class AwsS3MessageRepository(private val s3client: AmazonS3,
                             private val bucketName: String) : MessageRepository {

    override fun find(application: String, language: Locale): Messages {
        val allS3MessagesBundle = getAllS3MessagesBundleFor(application)
        val s3MessageBundle = s3MessageBundleFor(allS3MessagesBundle, application, language.toString())
        return s3MessageBundleContentFor(s3MessageBundle)
    }

    private fun s3MessageBundleFor(message: ObjectListing, application: String, language: String) =
            messageKeyFor(application)
                    .let { messagesKey ->
                        Optional.ofNullable(message.objectSummaries.find { it.key.equals("$messagesKey/messages_$language.properties") })
                                .orElse(message.objectSummaries.find { it.key.equals("$messagesKey/messages.properties") })
                    }

    private fun s3MessageBundleContentFor(message: S3ObjectSummary) =
            s3client.getObject(bucketName, message.key)
                    .objectContent.buffered()
                    .let { loadBundle(it) }

    private fun loadBundle(it: BufferedInputStream): Map<String, String> {
        val bundle = Properties()
        bundle.load(it)
        return bundle.toMap() as Map<String, String>
    }

    private fun messageKeyFor(application: String) = "i18n-messages/$application"

    private fun getAllS3MessagesBundleFor(application: String) = s3client.listObjects(bucketName, messageKeyFor(application))

}