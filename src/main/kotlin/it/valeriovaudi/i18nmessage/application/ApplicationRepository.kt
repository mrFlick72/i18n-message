package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO

interface ApplicationRepository {

    fun save(application: Application) : IO<Application>

    fun delete(application: Application) : IO<Application>

    fun findFor(id: String) : IO<Application>
}