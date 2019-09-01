package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import org.springframework.data.cassandra.core.CassandraTemplate

class CassandraApplicationRepository(private val cassandraTemplate: CassandraTemplate) : ApplicationRepository {

    val insert: (Application) -> String =
            { application ->
                "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) " +
                        "VALUES " +
                        "('${application.id}', '${application.name}', '${application.defaultLanguage.lang}')"
            }

    override fun save(application: Application): IO<Application> {
        return IO { cassandraTemplate.cqlOperations.execute(insert(application)) }
                .flatMap {
                    when (it) {
                        true -> IO { application }
                        false -> IO.raiseError(NotAppliedStatementException("The statement was not applied."))
                    }
                }
    }

    override fun delete(application: Application): IO<Application> {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun findFor(id: String): IO<Application> {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

}