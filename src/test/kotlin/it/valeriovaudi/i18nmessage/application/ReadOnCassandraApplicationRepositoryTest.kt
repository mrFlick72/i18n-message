package it.valeriovaudi.i18nmessage.application

import it.valeriovaudi.i18nmessage.Language.Companion.defaultLanguage
import junit.framework.Assert.assertTrue
import junit.framework.Assert.fail
import org.hamcrest.CoreMatchers.equalTo
import org.hamcrest.CoreMatchers.isA
import org.junit.Assert.assertThat
import org.junit.Before
import org.junit.ClassRule
import org.junit.Test
import org.junit.rules.ExpectedException
import org.junit.runner.RunWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.dao.EmptyResultDataAccessException
import org.springframework.data.cassandra.core.CassandraTemplate
import org.springframework.test.context.junit4.SpringRunner
import org.testcontainers.containers.DockerComposeContainer
import java.io.File

@SpringBootTest
@RunWith(SpringRunner::class)
class ReadOnCassandraApplicationRepositoryTest {

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
    }

    @Autowired
    lateinit var cassandraTemplate: CassandraTemplate

    @Test
    fun `find by id an Application`() {
        val cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)

        val expected = Application("AN_APPLICATION_ID", "AN_APPLICATION", defaultLanguage())

        val save = cassandraApplicationRepository.findFor("AN_APPLICATION_ID")
        save.attempt()
                .unsafeRunSync()
                .fold(
                        { it.printStackTrace(); fail() },
                        { assertThat(it, equalTo(expected)) }
                )
    }

    @Test
    fun `find by id an Application that does not exist`() {
        val cassandraApplicationRepository = CassandraApplicationRepository(cassandraTemplate)

        val save = cassandraApplicationRepository.findFor("A_NOT_EXISTED_APPLICATION")
        save.attempt()
                .unsafeRunSync()
                .fold(
                        {
                            it.printStackTrace()
                            assertTrue(it is EmptyResultDataAccessException)
                        },
                        { fail("I should do not find any application for the input") }
                )
    }
}
