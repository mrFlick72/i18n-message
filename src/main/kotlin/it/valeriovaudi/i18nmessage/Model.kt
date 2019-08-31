package it.valeriovaudi.i18nmessage

import arrow.core.Option
import java.util.*


data class Language(val lang: Locale) {
    companion object {
        fun availableLanguages() = listOf(Language(Locale.ITALIAN), Language(Locale.ENGLISH))
    }
}

data class MessageKey(val family: String, val key: String)

data class Application(val id: String, val name: String, val defaultLocale: Language)

data class Message(val application: Application, val language: Option<Language>, val key: MessageKey, val content: String)