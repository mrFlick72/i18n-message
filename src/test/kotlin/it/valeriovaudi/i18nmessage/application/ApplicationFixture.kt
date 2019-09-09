package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import com.datastax.driver.core.Row
import it.valeriovaudi.i18nmessage.languages.Language
import org.springframework.data.cassandra.core.cql.CqlOperations
import java.util.*

val findApplicationFor = { id: String, cqlOperations: CqlOperations ->
    IO {
        cqlOperations.queryForObject("SELECT * FROM i18n_messages.APPLICATION WHERE ID = '$id'")
        { row: Row, rowNum: Int -> Application(id = row.getString("id"), defaultLanguage = Language(Locale(row.getString("defaultLanguage"))), name = row.getString("name")) }
    }
}


