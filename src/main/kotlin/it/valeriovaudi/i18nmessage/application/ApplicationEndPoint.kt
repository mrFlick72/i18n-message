package it.valeriovaudi.i18nmessage.application

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RestController

@RestController
class ApplicationEndPoint(private val applicationRepository: CassandraApplicationRepository) {

    @PostMapping("/application")
    fun saveApplication(@RequestBody application: Application) =
            applicationRepository.save(application)
                    .attempt()
                    .unsafeRunSync()
                    .fold(handleInternalServerError(),
                            { ResponseEntity.status(HttpStatus.CREATED).build() })





    private fun handleInternalServerError(): (Throwable) -> ResponseEntity<Throwable> =
            { ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(it) }

}