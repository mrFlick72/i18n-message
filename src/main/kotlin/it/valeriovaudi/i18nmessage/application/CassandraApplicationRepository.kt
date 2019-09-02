package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import com.datastax.driver.core.Row
import it.valeriovaudi.i18nmessage.Language
import org.springframework.data.cassandra.core.CassandraTemplate
import java.util.*

private const val INSERT_QUERY: String = "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) VALUES (?,?,?)"
private const val SELECT_QUERY: String = "SELECT * FROM i18n_messages.APPLICATION WHERE id=?"

val applicationMapper = { row: Row, _: Int -> Application(id = row.getString("id"), defaultLanguage = Language(Locale(row.getString("defaultLanguage"))), name = row.getString("name")) }

class CassandraApplicationRepository(private val cassandraTemplate: CassandraTemplate) : ApplicationRepository {

    override fun save(application: Application): IO<Application> =
            IO { cassandraTemplate.cqlOperations.execute(INSERT_QUERY, application.id, application.name, application.defaultLanguage.asString()) }
                    .flatMap {
                        when (it) {
                            true -> IO { application }
                            false -> IO.raiseError(NotAppliedStatementException("The statement was not applied."))
                        }
                    }


    override fun delete(application: Application): IO<Application> =
            TODO("not implemented") //To change body of created functions use File | Settings | File Templates.


    override fun findFor(id: String): IO<Application> =
            IO { cassandraTemplate.cqlOperations.queryForObject(SELECT_QUERY, applicationMapper, arrayOf(id)) }

}