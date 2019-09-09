package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO

interface ApplicationRepository {

    fun save(application: Application): IO<Unit>

    fun deleteFor(id: String): IO<Unit>

    fun findFor(id: String): IO<Application>

    fun findAll(): IO<List<Application>>
}