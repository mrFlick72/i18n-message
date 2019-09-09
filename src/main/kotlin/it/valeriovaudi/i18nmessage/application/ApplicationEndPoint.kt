package it.valeriovaudi.i18nmessage.application

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*

@RestController
class ApplicationEndPoint(private val applicationRepository: CassandraApplicationRepository) {

    @GetMapping("/application")
    fun getApplications() =
            applicationRepository.findAll()
                    .unsafeRunSync()
                    .let { ResponseEntity.ok(it) }

    @GetMapping("/application/{id}")
    fun getApplication(@PathVariable id: String) =
            applicationRepository.findFor(id)
                    .unsafeRunSync()
                    .let { ResponseEntity.ok(it) }


    @PutMapping("/application")
    fun saveApplication(@RequestBody application: Application) =
            applicationRepository.save(application)
                    .let { ResponseEntity.status(HttpStatus.CREATED).build<Unit>() }


    @DeleteMapping("/application/{id}")
    fun saveApplication(@PathVariable id: String) =
            applicationRepository.deleteFor(id)
                    .unsafeRunSync()
                    .let { ResponseEntity.noContent().build<Unit>() }

}