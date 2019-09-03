package it.valeriovaudi.i18nmessage.application

import arrow.core.orNull
import it.valeriovaudi.i18nmessage.Language.Companion.defaultLanguage
import junit.framework.Assert.fail
import org.hamcrest.CoreMatchers.equalTo
import org.junit.Assert.assertThat
import org.junit.Before
import org.junit.ClassRule
import org.junit.Test
import org.junit.runner.RunWith
import org.mockito.BDDMockito.given
import org.mockito.Mock
import org.mockito.Mockito.verify
import org.mockito.junit.MockitoJUnitRunner
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.cassandra.core.CassandraTemplate
import org.springframework.data.cassandra.core.cql.CqlOperations
import org.springframework.test.annotation.DirtiesContext
import org.springframework.test.context.junit4.SpringRunner
import org.testcontainers.containers.DockerComposeContainer
import java.io.File

@SpringBootTest
@DirtiesContext
@RunWith(SpringRunner::class)
class InsertOnCassandraApplicationRepositoryIT {

    companion object {
        @ClassRule
        @JvmField
        val container: DockerComposeContainer<*> = DockerComposeContainer<Nothing>(File("src/test/resources/cassandra/docker-compose.yml"))
                .withExposedService("cassandra_1", 9042)
    }

    @Before
    fun setUp() {
        cassandraTemplate.cqlOperations.execute("CREATE KEYSPACE IF NOT EXISTS i18n_messages WITH replication = {'class':'SimpleStrategy','replication_factor':'1'};")
        cassandraTemplate.cqlOperations.execute("CREATE TABLE IF NOT EXISTS i18n_messages.APPLICATION (id varchar primary key, name varchar,defaultLanguage varchar);")
    }

    @Autowired
    lateinit var cassandraTemplate: CassandraTemplate

    @Test
    fun `save a new Application`() {
        val cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)

        val id = "AN_APPLICATION_ID"
        val expected = Application(id, "AN_APPLICATION", defaultLanguage())

        val save = cassandraApplicationRepository.save(expected)
        save.attempt()
                .unsafeRunSync()
                .fold(
                        { it.printStackTrace(); fail() },
                        { assertThat(it, equalTo(expected)) }
                )
        val actual = findApplicationFor(id, cassandraTemplate.cqlOperations).attempt().unsafeRunSync().orNull()
        assertThat(actual, equalTo(expected))
    }
}


@RunWith(MockitoJUnitRunner::class)
class InsertOnCassandraApplicationRepositoryTest {

    @Mock
    lateinit var cassandraTemplate: CassandraTemplate

    @Mock
    lateinit var cqlOperations: CqlOperations

    @Test
    fun `save a new Application`() {
        val cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)

        val expected = Application("AN_APPLICATION_ID", "AN_APPLICATION", defaultLanguage())

        val insertQuery = "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) VALUES (?,?,?)"
        given(cassandraTemplate.cqlOperations).willReturn(cqlOperations)
        given(cqlOperations.execute(insertQuery, "AN_APPLICATION_ID", "AN_APPLICATION", "en")).willReturn(false)

        val save = cassandraApplicationRepository.save(expected)
        save.attempt()
                .unsafeRunSync()
                .fold(
                        {},
                        { fail() }

                )
        verify(cqlOperations).execute(insertQuery, "AN_APPLICATION_ID", "AN_APPLICATION", "en")
    }
}