package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import arrow.syntax.function.partially2
import junit.framework.Assert.*
import org.junit.Before
import org.junit.ClassRule
import org.junit.Test
import org.junit.rules.ExpectedException
import org.junit.runner.RunWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.dao.EmptyResultDataAccessException
import org.springframework.data.cassandra.core.CassandraTemplate
import org.springframework.data.cassandra.core.cql.CqlOperations
import org.springframework.test.annotation.DirtiesContext
import org.springframework.test.context.junit4.SpringRunner
import org.testcontainers.containers.DockerComposeContainer
import java.io.File

@SpringBootTest
@DirtiesContext
@RunWith(SpringRunner::class)
class DeleteOnCassandraApplicationRepositoryTest {

    companion object {
        @ClassRule
        @JvmField
        val container: DockerComposeContainer<*> = DockerComposeContainer<Nothing>(File("src/test/resources/cassandra/docker-compose.yml"))
                .withExposedService("cassandra_1", 9042)

        @ClassRule
        @JvmField
        var exception = ExpectedException.none()
    }


    @Before
    fun setUp() {
        cassandraTemplate.cqlOperations.execute("CREATE KEYSPACE IF NOT EXISTS i18n_messages WITH replication = {'class':'SimpleStrategy','replication_factor':'1'};")
        cassandraTemplate.cqlOperations.execute("CREATE TABLE IF NOT EXISTS i18n_messages.APPLICATION (id varchar primary key, name varchar,defaultLanguage varchar);")
        cassandraTemplate.cqlOperations.execute("INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) VALUES ('AN_APPLICATION_ID', 'AN_APPLICATION', 'en');")

        cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)
        findApplication = findApplicationFor.partially2(cassandraTemplate.cqlOperations)
    }

    @Autowired
    lateinit var cassandraTemplate: CassandraTemplate

    lateinit var cassandraApplicationRepository: CassandraApplicationRepository
    lateinit var findApplication: (id: String) -> IO<Application>

    val happyPathAssertion: (Throwable) -> Unit = { it.printStackTrace(); assertTrue(it is EmptyResultDataAccessException) }
    val failurePathAssertion: (Application) -> Unit = { fail() }

    @Test
    fun `delete by id an Application`() {
        val cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)
        val findApplicationFor = findApplicationFor.partially2(cassandraTemplate.cqlOperations)

        val id = "AN_APPLICATION_ID"

        cassandraApplicationRepository
                .deleteFor(id)
                .flatMap { findApplicationFor(id) }
                .attempt()
                .unsafeRunSync()
                .fold(happyPathAssertion, failurePathAssertion)
    }

    @Test
    fun `delete by id an Application that does not exist`() {
        val cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)
        val findApplicationFor = findApplicationFor.partially2(cassandraTemplate.cqlOperations)

        val id = "AN_OTHER_APPLICATION_ID"

        cassandraApplicationRepository
                .deleteFor(id)
                .flatMap { findApplicationFor(id) }
                .attempt()
                .unsafeRunSync()
                .fold(happyPathAssertion, failurePathAssertion)
    }

}