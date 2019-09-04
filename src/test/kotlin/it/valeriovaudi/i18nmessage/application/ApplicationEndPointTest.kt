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
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post
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


        mockMvc.perform(post("/application")
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


        mockMvc.perform(post("/application")
                .contentType("application/json")
                .content(objectMapper.writeValueAsString(application)))
                .andExpect(status().isInternalServerError)
    }
}