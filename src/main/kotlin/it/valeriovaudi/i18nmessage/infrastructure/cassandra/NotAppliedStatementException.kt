package it.valeriovaudi.i18nmessage.infrastructure.cassandra

import java.lang.RuntimeException

data class NotAppliedStatementException(override val message : String) : RuntimeException(message)