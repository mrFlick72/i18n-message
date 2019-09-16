package it.valeriovaudi.i18nmessage.messages


interface MessageRepository {

    fun findOne(application: String, family: String, key: String): Map<String, String>

    fun find(application: String, family: String): Map<String, String>
}