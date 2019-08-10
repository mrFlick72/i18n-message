package it.valeriovaudi.i18nmessage

import arrow.core.Option
import java.util.*


data class Language(val lang: Locale) {
    companion object {
        fun availableLanguages() = listOf(Language(Locale.ITALIAN), Language(Locale.ENGLISH))
    }
}

data class Application(val id: String, val name: String, val defaultLocale: Language)

data class Message(val application: Application, val language: Option<Language>, val key: String, val content: String)