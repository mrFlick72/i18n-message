package it.valeriovaudi.i18nmessage.messages

import org.springframework.jms.annotation.JmsListener
import org.springframework.stereotype.Component

@Component
class AwsSQSListener {

    @JmsListener(destination = "i18n-messages-updates")
    fun onMessage(message: String) {
        println(message)
    }
}