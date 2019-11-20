FROM openjdk:11

ADD i18n-messages.jar /usr/local/i18n-messages/

WORKDIR /usr/local/i18n-messages/

CMD ["java", "-Xms512m", "-Xmx1g", "-jar", "i18n-messages.jar"]