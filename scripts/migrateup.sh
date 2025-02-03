#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema
goose up $DATABASE_URL