#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema
goose $DATABASE_URL up