package it.valeriovaudi.i18nmessage.application

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*

@RestController
class ApplicationEndPoint(private val applicationRepository: CassandraApplicationRepository) {

    @GetMapping("/application")
    fun getApplications() =
            applicationRepository.findAll()
                    .attempt()
                    .unsafeRunSync()
                    .fold(handleInternalServerError(),
                            { ResponseEntity.ok(it) })

    @GetMapping("/application/{id}")
    fun getApplication(@PathVariable id: String) =
            applicationRepository.findFor(id)
                    .attempt()
                    .unsafeRunSync()
                    .fold(handleInternalServerError(),
                            { ResponseEntity.ok(it) })


    @PutMapping("/application")
    fun saveApplication(@RequestBody application: Application) =
            applicationRepository.save(application)
                    .attempt()
                    .unsafeRunSync()
                    .fold(handleInternalServerError(),
                            { ResponseEntity.status(HttpStatus.CREATED).build() })


    @DeleteMapping("/application/{id}")
    fun saveApplication(@PathVariable id: String) =
            applicationRepository.deleteFor(id)
                    .attempt()
                    .unsafeRunSync()
                    .fold(handleInternalServerError(),
                            { ResponseEntity.noContent().build() })


    private fun handleInternalServerError(): (Throwable) -> ResponseEntity<Throwable> =
            { ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(it) }

}