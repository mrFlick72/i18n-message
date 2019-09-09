package it.valeriovaudi.i18nmessage.languages

import org.hamcrest.CoreMatchers.equalTo
import org.junit.Assert.*
import org.junit.Test
import org.springframework.security.test.context.support.WithMockUser
import java.util.*


class YamlLanguageRepositoryTest {

    @Test
    @WithMockUser
    fun `get all languages`() {
        val yamlLanguageModel = YamlLanguageModel()
        yamlLanguageModel.lang = listOf("it", "en")

        val yamlLanguageRepository = YamlLanguageRepository(yamlLanguageModel)

        val actual = yamlLanguageRepository.findAll().unsafeRunSync()
        assertThat(actual, equalTo(listOf(Language(Locale.ITALIAN), Language(Locale.ENGLISH))))

    }
}