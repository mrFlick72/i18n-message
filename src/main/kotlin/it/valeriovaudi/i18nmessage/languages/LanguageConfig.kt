package it.valeriovaudi.i18nmessage.languages

import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.cloud.context.config.annotation.RefreshScope
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
@EnableConfigurationProperties(YamlLanguageModel::class)
class LanguageConfig {

    @Bean
    @RefreshScope
    fun languageRepository(yamlLanguageModel: YamlLanguageModel) =
            YamlLanguageRepository(yamlLanguageModel)
}