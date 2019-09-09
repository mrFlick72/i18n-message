package it.valeriovaudi.i18nmessage.languages

import java.util.*

data class Language(val lang: Locale) {
    companion object {
        fun defaultLanguage() = Language(Locale.ENGLISH)
    }

    fun asString(): String = lang.toString()
}