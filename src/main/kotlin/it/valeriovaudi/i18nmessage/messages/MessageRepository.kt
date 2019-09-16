package it.valeriovaudi.i18nmessage.messages

import arrow.effects.IO

interface MessageRepository {

    fun save(message: Message)

    fun delete(messageKey: MessageKey)

    fun findOne(messageKey: MessageKey): IO<Message>

    fun find(family: String): IO<List<Message>>
}