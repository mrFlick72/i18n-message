package repository

import (
	"bytes"
	fmt "fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/magiconair/properties"
)

var (
	patternWithLanguage    = "%s/documents/%s?path=%s&fileName=messages_%v&fileExt=properties"
	patternWithoutLanguage = "%s/documents/%s?path=%s&fileName=messages&fileExt=properties"
	acceptableLength       = 2
	propertyKeyPosition    = 0
	propertyValuePosition  = 1
	translationNotFound    = 404
	noLanguage             = ""
	spliteratorCharacter   = "="
)

type Language = string
type Message = map[string]string

type MessageRepository interface {
	Find(application string, language string, context map[string]string) (*Message, error)
}

type S3MessageRepository struct {
	Client    *s3.S3
	BuketName string
}

func (repository *S3MessageRepository) Find(application string, language string, context map[string]string) (*Message, error) {
	if language != "" {
		language = "_" + language
	}
	objectKey := objectKeyFor(application, language)

	object, err := downloadFile(repository, objectKey)

	if err != nil {
		fallBackObjectKey := objectKeyFor(application, "")
		object, err = downloadFile(repository, fallBackObjectKey)

		if err != nil {
			return nil, err
		}
	}

	bytes := objectContentFor(object)

	properties, err := properties.Load(bytes, properties.UTF8)

	if err != nil {
		return nil, err
	}

	message := messageFrom(properties)
	return &message, nil
}

func messageFrom(properties *properties.Properties) Message {
	var message = make(Message)
	message = properties.Map()
	return message
}

func objectContentFor(object *s3.GetObjectOutput) []byte {
	buf := new(bytes.Buffer)

	body := object.Body
	buf.ReadFrom(body)
	defer body.Close()

	bytes := buf.Bytes()
	return bytes
}

func objectKeyFor(application string, language string) string {
	return fmt.Sprintf("%s/messages%s.properties", application, language)
}

func downloadFile(repository *S3MessageRepository, objectKey string) (*s3.GetObjectOutput, error) {
	object, err := repository.Client.GetObject(&s3.GetObjectInput{
		Bucket: &repository.BuketName,
		Key:    &objectKey,
	})

	if err != nil {
		return nil, err
	}

	return object, nil
}
