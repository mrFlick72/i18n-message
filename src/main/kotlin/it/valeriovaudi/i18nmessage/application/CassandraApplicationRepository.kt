package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import com.datastax.driver.core.Row
import it.valeriovaudi.i18nmessage.languages.Language
import org.slf4j.LoggerFactory
import org.springframework.data.cassandra.core.CassandraTemplate
import java.util.*

private const val SELECT_ALL_QUERY: String = "SELECT * FROM i18n_messages.APPLICATION"
private const val SELECT_QUERY: String = "SELECT * FROM i18n_messages.APPLICATION WHERE id=?"
private const val INSERT_QUERY: String = "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) VALUES (?,?,?)"
private const val DELETE_QUERY: String = "DELETE  FROM i18n_messages.APPLICATION WHERE id=?"

val idFor = { row: Row -> row.getString("id") }
val defaultLanguageFor = { row: Row -> Language(Locale(row.getString("defaultLanguage"))) }
val nameFor = { row: Row -> row.getString("name") }

val applicationMapper = { row: Row, _: Int -> Application(id = idFor(row), defaultLanguage = defaultLanguageFor(row), name = nameFor(row)) }

private val LOGGER = LoggerFactory.getLogger(CassandraApplicationRepository::class.java);

open class CassandraApplicationRepository(private val cassandraTemplate: CassandraTemplate) : ApplicationRepository {

    override fun findAll(): IO<List<Application>> =
            IO { cassandraTemplate.cqlOperations.query(SELECT_ALL_QUERY, applicationMapper) }

    override fun save(application: Application): IO<Unit> =
            IO { cassandraTemplate.cqlOperations.execute(INSERT_QUERY, application.id, application.name, application.defaultLanguage.asString()) }
                    .map { executed ->
                        LOGGER.debug("the query is executed: $executed")
                    }

    override fun deleteFor(id: String): IO<Unit> =
            IO { cassandraTemplate.cqlOperations.execute(DELETE_QUERY, id) }
                    .map { executed ->
                        LOGGER.debug("the query is executed: $executed")
                    }

    override fun findFor(id: String): IO<Application> =
            IO { cassandraTemplate.cqlOperations.queryForObject(SELECT_QUERY, applicationMapper, arrayOf(id)) }
}