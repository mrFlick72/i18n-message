package it.valeriovaudi.i18nmessage.languages

import arrow.effects.IO
import org.springframework.boot.context.properties.ConfigurationProperties
import java.util.*

open class YamlLanguageRepository(private val yamlLanguageModel: YamlLanguageModel) : LanguageRepository {
    override fun findAll(): IO<List<Language>> =
            IO {
                yamlLanguageModel.lang
                        .map { Locale(it) }
                        .map { Language(it) }
            }

}

@ConfigurationProperties("languages")
open class YamlLanguageModel {
    lateinit var lang: List<String>
}