package it.valeriovaudi.i18nmessage.messages

import arrow.core.Option
import it.valeriovaudi.i18nmessage.application.Application
import it.valeriovaudi.i18nmessage.languages.Language

data class MessageKey(val family: String, val key: String)

data class Message(val application: Application, val language: Option<Language>, val key: MessageKey, val content: String)