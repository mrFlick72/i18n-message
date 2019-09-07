package it.valeriovaudi.i18nmessage.application

import arrow.effects.IO
import it.valeriovaudi.i18nmessage.Language
import org.junit.Test
import org.junit.runner.RunWith
import org.mockito.BDDMockito.given
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest
import org.springframework.boot.test.mock.mockito.MockBean
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
    fun `we are able to save a new Application`() {
        val application = Application("AN_ID", "AN_APPLICATION", Language.defaultLanguage())
        val io = IO { application }

        given(applicationRepository.save(application))
                .willReturn(io)

        mockMvc.perform(put("/application")
                .contentType("application/json")
                .content(objectMapper.writeValueAsString(application)))
                .andExpect(status().isCreated)
    }

    @Test
    fun `save a new Application goes in error`() {
        val application = Application("AN_ID", "AN_APPLICATION", Language.defaultLanguage())
        val io = IO.raiseError<Application>(RuntimeException())

        given(applicationRepository.save(application))
                .willReturn(io)

        mockMvc.perform(put("/application")
                .contentType("application/json")
                .content(objectMapper.writeValueAsString(application)))
                .andExpect(status().isInternalServerError)
    }

    @Test
    fun `delete an Application`() {
        given(applicationRepository.deleteFor("AN_ID"))
                .willReturn(IO { Unit })

        mockMvc.perform(delete("/application/AN_ID"))
                .andExpect(status().isNoContent)
    }

    @Test
    fun `get an Application`() {
        val application = Application("AN_ID", "AN_APPLICATION", Language.defaultLanguage())

        given(applicationRepository.findFor("AN_ID"))
                .willReturn(IO { application })

        mockMvc.perform(get("/application/AN_ID"))
                .andExpect(status().isOk)
                .andExpect(content().string(objectMapper.writeValueAsString(application)))
    }

    @Test
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