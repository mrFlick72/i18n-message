package it.valeriovaudi.i18nmessage

import com.amazon.sqs.javamessaging.ProviderConfiguration
import com.amazon.sqs.javamessaging.SQSConnectionFactory
import com.amazonaws.auth.AWSStaticCredentialsProvider
import com.amazonaws.auth.BasicAWSCredentials
import com.amazonaws.services.sqs.AmazonSQSClientBuilder
import it.valeriovaudi.i18nmessage.messages.MessageRepository
import it.valeriovaudi.i18nmessage.messages.RestMessageRepository
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.ConstructorBinding
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean
import org.springframework.jms.annotation.EnableJms
import org.springframework.jms.config.DefaultJmsListenerContainerFactory
import org.springframework.jms.support.destination.DynamicDestinationResolver
import org.springframework.messaging.rsocket.RSocketRequester
import org.springframework.messaging.rsocket.RSocketStrategies
import org.springframework.web.reactive.function.client.WebClient
import reactor.core.publisher.Mono
import javax.jms.ConnectionFactory
import javax.jms.Session


@EnableJms
@SpringBootApplication
@EnableConfigurationProperties(value = [RSocketApplicationClientApps::class])
class I18nMessageApplication {

    @Bean
    fun sqsConnectionFactory(@Value("\${aws.access-key}") accessKey: String,
                             @Value("\${aws.secret-key}") awsSecretKey: String,
                             @Value("\${aws.region}") awsRegion: String): SQSConnectionFactory {
        return SQSConnectionFactory(
                ProviderConfiguration(),
                AmazonSQSClientBuilder
                        .standard()
                        .withCredentials(
                                AWSStaticCredentialsProvider(
                                        BasicAWSCredentials(
                                                accessKey,
                                                awsSecretKey
                                        )
                                )
                        )
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
    fun requesters(rSocketStrategies: RSocketStrategies,
                   rSocketApplicationClientApps: RSocketApplicationClientApps,
                   builder: RSocketRequester.Builder): Map<String, Mono<RSocketRequester>> =
            emptyMap()
    /*         rSocketApplicationClientApps.clients
                     .map {
                         val address = InetSocketAddress(it.host, it.port)
                         val clientTransport: TcpClientTransport = TcpClientTransport.create(address)
                         mapOf(it.id to builder.rsocketStrategies(rSocketStrategies)
                                 .connect(clientTransport))
                     }.reduce { acc, map -> acc + map }*/

    @Bean
    fun messageRepository(@Value("\${repository-service.baseUrl}") repositoryServiceUrl: String,
                          @Value("\${repository-service.serviceRegistrationName}") registrationName: String): MessageRepository {

        return RestMessageRepository(repositoryServiceUrl, registrationName, WebClient.builder());
    }

}

@ConstructorBinding
@ConfigurationProperties(prefix = "rsocket")
data class RSocketApplicationClientApps(val clients: List<RSocketApplicationClientApp> = emptyList())

data class RSocketApplicationClientApp(val id: String, val host: String, val port: Int)

fun main(args: Array<String>) {
    runApplication<I18nMessageApplication>(*args)
}
