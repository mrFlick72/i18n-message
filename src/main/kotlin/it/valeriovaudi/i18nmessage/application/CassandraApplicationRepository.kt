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
}