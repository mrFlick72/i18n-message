#!/bin/bash

aws s3api create-bucket --bucket bucket --endpoint http://localhost:4566 --region us-east-1

aws s3 cp  messages.properties s3://bucket/AN_APPLICATION/messages.properties --endpoint http://localhost:4566 --region us-east-1
aws s3 cp  messages_it.properties s3://bucket/AN_APPLICATION/messages_it.properties --endpoint http://localhost:4566 --region us-east-1

aws s3 ls s3://bucket --endpoint http://localhost:4566 --region us-east-1
