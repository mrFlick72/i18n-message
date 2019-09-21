package it.valeriovaudi.i18nmessage.messages

import java.util.*

typealias Messages = Map<String, String>

interface MessageRepository {

    fun find(application: String, language: Locale): Messages
}