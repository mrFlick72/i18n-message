package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import com.datastax.driver.core.Row
import it.valeriovaudi.i18nmessage.Language
import org.springframework.data.cassandra.core.CassandraTemplate
import java.util.*

typealias InsertQueryGenerator = (Application) -> String
typealias SelectQueryGenerator = (String) -> String

val insertQueryFor: InsertQueryGenerator =
        { application ->
            "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) " +
                    "VALUES " +
                    "('${application.id}', '${application.name}', '${application.defaultLanguage.lang}')"
        }
val selectQueryFor: SelectQueryGenerator =
        { id -> "SELECT * FROM i18n_messages.APPLICATION WHERE id=?" }

val applicationMapper = { row: Row, _: Int -> Application(id = row.getString("id"), defaultLanguage = Language(Locale(row.getString("defaultLanguage"))), name = row.getString("name")) }

class CassandraApplicationRepository(private val cassandraTemplate: CassandraTemplate) : ApplicationRepository {

    override fun save(application: Application): IO<Application> =
            IO { cassandraTemplate.cqlOperations.execute(insertQueryFor(application)) }
                    .flatMap {
                        when (it) {
                            true -> IO { application }
                            false -> IO.raiseError(NotAppliedStatementException("The statement was not applied."))
                        }
                    }


    override fun delete(application: Application): IO<Application> =
            TODO("not implemented") //To change body of created functions use File | Settings | File Templates.


    override fun findFor(id: String): IO<Application> =
            IO { cassandraTemplate.cqlOperations.queryForObject("SELECT * FROM i18n_messages.APPLICATION WHERE id=?", applicationMapper, arrayOf(id)) }

}