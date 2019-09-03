package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import com.datastax.driver.core.Row
import it.valeriovaudi.i18nmessage.Language
import org.springframework.data.cassandra.core.CassandraTemplate
import java.util.*

private const val INSERT_QUERY: String = "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) VALUES (?,?,?)"
private const val SELECT_QUERY: String = "SELECT * FROM i18n_messages.APPLICATION WHERE id=?"
private const val DELETE_QUERY: String = "DELETE  FROM i18n_messages.APPLICATION WHERE id=?"

val idFor = { row: Row -> row.getString("id") }
val defaultLanguageFor = { row: Row -> Language(Locale(row.getString("defaultLanguage"))) }
val nameFor = { row: Row -> row.getString("name") }

val applicationMapper = { row: Row, _: Int -> Application(id = idFor(row), defaultLanguage = defaultLanguageFor(row), name = nameFor(row)) }

class CassandraApplicationRepository(private val cassandraTemplate: CassandraTemplate) : ApplicationRepository {

    override fun save(application: Application): IO<Application> =
            IO { cassandraTemplate.cqlOperations.execute(INSERT_QUERY, application.id, application.name, application.defaultLanguage.asString()) }
                    .flatMap {
                        when (it) {
                            true -> IO { application }
                            false -> IO.raiseError(NotAppliedStatementException("The statement was not applied."))
                        }
                    }


    override fun deleteFor(id: String): IO<Unit> =
            IO { cassandraTemplate.cqlOperations.execute(DELETE_QUERY, id) }
                    .flatMap {
                        when (it) {
                            true -> IO { Unit }
                            false -> IO.raiseError(NotAppliedStatementException("The statement was not applied."))
                        }
                    }

    override fun findFor(id: String): IO<Application> =
            IO { cassandraTemplate.cqlOperations.queryForObject(SELECT_QUERY, applicationMapper, arrayOf(id)) }

}