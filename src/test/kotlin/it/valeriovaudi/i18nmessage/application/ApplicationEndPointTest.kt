package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import it.valeriovaudi.i18nmessage.languages.Language
import org.junit.Test
import org.junit.runner.RunWith
import org.mockito.BDDMockito.given
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest
import org.springframework.boot.test.mock.mockito.MockBean
import org.springframework.security.test.context.support.WithMockUser
import org.springframework.test.context.junit4.SpringRunner
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.content
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.status
import org.testcontainers.shaded.com.fasterxml.jackson.databind.ObjectMapper

@WebMvcTest
@RunWith(SpringRunner::class)
class ApplicationEndPointTest {

    @Autowired
    lateinit var mockMvc: MockMvc

    @MockBean
    lateinit var applicationRepository: CassandraApplicationRepository

    val objectMapper = ObjectMapper()

    @Test
    @WithMockUser
    fun `we are able to save a new Application`() {
        val application = Application("AN_ID", "AN_APPLICATION", Language.defaultLanguage())

        given(applicationRepository.save(application))
                .willReturn(IO { Unit })

        mockMvc.perform(put("/application")
                .contentType("application/json")
                .content(objectMapper.writeValueAsString(application)))
                .andExpect(status().isCreated)
    }

    @Test
    @WithMockUser
    fun `delete an Application`() {
        given(applicationRepository.deleteFor("AN_ID"))
                .willReturn(IO { Unit })

        mockMvc.perform(delete("/application/AN_ID"))
                .andExpect(status().isNoContent)
    }

    @Test
    @WithMockUser
    fun `get an Application`() {
        val application = Application("AN_ID", "AN_APPLICATION", Language.defaultLanguage())

        given(applicationRepository.findFor("AN_ID"))
                .willReturn(IO { application })

        mockMvc.perform(get("/application/AN_ID"))
                .andExpect(status().isOk)
                .andExpect(content().string(objectMapper.writeValueAsString(application)))
    }

    @Test
    @WithMockUser
    fun `get all Applications`() {
        val applications = listOf(
                Application("AN_ID", "AN_APPLICATION", Language.defaultLanguage()),
                Application("ANOTHER_ID", "ANOTHER_APPLICATION", Language.defaultLanguage())
        )

        given(applicationRepository.findAll())
                .willReturn(IO { applications })

        mockMvc.perform(get("/application"))
                .andExpect(status().isOk)
                .andExpect(content().string(objectMapper.writeValueAsString(applications)))
    }
}