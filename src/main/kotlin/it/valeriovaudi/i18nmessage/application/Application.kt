package it.valeriovaudi.i18nmessage.application

import it.valeriovaudi.i18nmessage.languages.Language

data class Application(val id: String, val name: String, val defaultLanguage: Language)
