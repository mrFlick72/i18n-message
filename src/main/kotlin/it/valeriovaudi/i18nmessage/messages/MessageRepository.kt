package it.valeriovaudi.i18nmessage.messages

import reactor.core.publisher.Mono
import java.util.*

typealias Messages = Map<String, String>

interface MessageRepository {

    fun find(application: String, language: Locale): Mono<Messages>
}