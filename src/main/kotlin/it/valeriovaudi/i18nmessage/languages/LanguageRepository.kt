package it.valeriovaudi.i18nmessage.languages

import arrow.effects.IO

interface LanguageRepository {

    fun findAll(): IO<List<Language>>

}