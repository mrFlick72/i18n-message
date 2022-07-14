#!/bin/bash


aws s3 rm  s3://bucket/AN_APPLICATION/messages.properties  --endpoint http://localhost:4566
aws s3 rm  s3://bucket/AN_APPLICATION/messages_it.properties --endpoint http://localhost:4566
aws s3 rm  s3://bucket/AN_APPLICATION --endpoint http://localhost:4566

aws s3api delete-bucket --bucket bucket --endpoint http://localhost:4566
