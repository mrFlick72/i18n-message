package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import org.springframework.data.cassandra.core.CassandraTemplate

class CassandraApplicationRepository(private val cassandraTemplate: CassandraTemplate) : ApplicationRepository {
    override fun save(application: Application): IO<Application> {
        val insert = "INSERT INTO i18n_messages.APPLICATION (id, name, defaultLanguage) VALUES ('${application.id}', '${application.name}', '${application.defaultLanguage.lang}')"
        return IO { cassandraTemplate.cqlOperations.execute(insert) }.map { println("esito $it"); application }
    }
}