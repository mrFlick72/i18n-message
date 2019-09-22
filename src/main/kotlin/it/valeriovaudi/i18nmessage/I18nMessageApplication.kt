package it.valeriovaudi.i18nmessage

import com.amazonaws.auth.AWSStaticCredentialsProvider
import com.amazonaws.auth.BasicAWSCredentials
import com.amazonaws.services.s3.AmazonS3ClientBuilder
import it.valeriovaudi.i18nmessage.messages.AwsS3MessageRepository
import it.valeriovaudi.i18nmessage.messages.MessageRepository
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean

@SpringBootApplication
class I18nMessageApplication {

    @Bean
    fun messageRepository(@Value("\${aws.s3.access-key}") accessKey: String,
                          @Value("\${aws.s3.secret-key}") awsSecretKey: String,
                          @Value("\${aws.s3.region}") awsRegion: String,
                          @Value("\${aws.s3.bucket}") awsBucket: String): MessageRepository {
        val credentials = BasicAWSCredentials(accessKey, awsSecretKey)

        val s3client = AmazonS3ClientBuilder
                .standard()
                .withCredentials(AWSStaticCredentialsProvider(credentials))
                .withRegion(awsRegion)
                .build()

        return AwsS3MessageRepository(s3client, awsBucket)
    }
}

fun main(args: Array<String>) {
    runApplication<I18nMessageApplication>(*args)
}
