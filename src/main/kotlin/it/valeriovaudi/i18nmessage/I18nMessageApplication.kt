package it.valeriovaudi.i18nmessage

import com.amazon.sqs.javamessaging.ProviderConfiguration
import com.amazon.sqs.javamessaging.SQSConnectionFactory
import com.amazonaws.auth.AWSCredentialsProvider
import com.amazonaws.auth.AWSStaticCredentialsProvider
import com.amazonaws.auth.BasicAWSCredentials
import com.amazonaws.services.s3.AmazonS3ClientBuilder
import com.amazonaws.services.sqs.AmazonSQSClientBuilder
import io.rsocket.transport.netty.client.TcpClientTransport
import it.valeriovaudi.i18nmessage.messages.AwsS3MessageRepository
import it.valeriovaudi.i18nmessage.messages.MessageRepository
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean
import org.springframework.jms.annotation.EnableJms
import org.springframework.jms.config.DefaultJmsListenerContainerFactory
import org.springframework.jms.support.destination.DynamicDestinationResolver
import org.springframework.messaging.rsocket.RSocketRequester
import org.springframework.messaging.rsocket.RSocketStrategies
import reactor.core.publisher.Mono
import java.net.InetSocketAddress
import javax.jms.ConnectionFactory
import javax.jms.Session


@EnableJms
@SpringBootApplication
class I18nMessageApplication {

    @Bean
    fun sqsConnectionFactory(@Value("\${aws.s3.region}") awsRegion: String,
                             awsCredentialsProvider: AWSCredentialsProvider): SQSConnectionFactory {
        return SQSConnectionFactory(
                ProviderConfiguration(),
                AmazonSQSClientBuilder
                        .standard()
                        .withCredentials(awsCredentialsProvider)
                        .withRegion(awsRegion)
                        .build()
        )

    }

    @Bean
    fun jmsListenerContainerFactory(sqsConnectionFactory: ConnectionFactory): DefaultJmsListenerContainerFactory {
        val factory = DefaultJmsListenerContainerFactory()
        factory.setConnectionFactory(sqsConnectionFactory)
        factory.setDestinationResolver(DynamicDestinationResolver())
        factory.setConcurrency("3-10")
        factory.setSessionAcknowledgeMode(Session.CLIENT_ACKNOWLEDGE)
        return factory
    }

    @Bean
    fun requester(rSocketStrategies: RSocketStrategies,
           @Value("\${i18n-messages.rsocket.host}") i18nHost: String,
           @Value("\${i18n-messages.rsocket.port}") i18nPort: Int,
           builder: RSocketRequester.Builder) : Mono<RSocketRequester> {
        val address = InetSocketAddress(i18nHost, i18nPort)
        val clientTransport: TcpClientTransport = TcpClientTransport.create(address)
        return  builder.rsocketStrategies(rSocketStrategies)
                .connect(clientTransport)

    }

    @Bean
    fun messageRepository(@Value("\${aws.s3.region}") awsRegion: String,
                          @Value("\${aws.s3.bucket}") awsBucket: String,
                          awsCredentialsProvider: AWSCredentialsProvider): MessageRepository {

        val s3client = AmazonS3ClientBuilder
                .standard()
                .withCredentials(awsCredentialsProvider)
                .withRegion(awsRegion)
                .build()

        return AwsS3MessageRepository(s3client, awsBucket)
    }

    @Bean
    fun awsCredentialsProvider(@Value("\${aws.s3.access-key}") accessKey: String,
                               @Value("\${aws.s3.secret-key}") awsSecretKey: String):
            AWSCredentialsProvider = AWSStaticCredentialsProvider(BasicAWSCredentials(accessKey, awsSecretKey))

}

fun main(args: Array<String>) {
    runApplication<I18nMessageApplication>(*args)
}
