package it.valeriovaudi.i18nmessage.application

import java.lang.RuntimeException

data class NotAppliedStatementException(override val message : String) : RuntimeException(message)