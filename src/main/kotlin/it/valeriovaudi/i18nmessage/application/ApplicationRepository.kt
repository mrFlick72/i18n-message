package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO

interface ApplicationRepository {

    fun save(application: Application) : IO<Application>
}