package it.valeriovaudi.i18nmessage.messages

import arrow.effects.IO

class CassandraMessageRepository : MessageRepository {
    override fun save(message: Message) {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun delete(messageKey: MessageKey) {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun findOne(messageKey: MessageKey): IO<Message> {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun find(family: String): IO<List<Message>> {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }
}