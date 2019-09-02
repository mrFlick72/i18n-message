package it.valeriovaudi.i18nmessage

import arrow.core.Option
import it.valeriovaudi.i18nmessage.application.Application
import java.util.*


data class Language(val lang: Locale) {
    companion object {
        fun defaultLanguage() = Language(Locale.ENGLISH)
        fun availableLanguages() = listOf(Language(Locale.ITALIAN), Language(Locale.ENGLISH))
    }

    fun asString(): String = lang.toString()
}

data class MessageKey(val family: String, val key: String)

data class Message(val application: Application, val language: Option<Language>, val key: MessageKey, val content: String)