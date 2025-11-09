#!/bin/bash

echo "Initializing database..."
until pg_isready -h db -U postgres -d gobid; do
    sleep 1
done

echo "Creating database if not exists"
psql -h db -U postgres -d gobid -c "CREATE DATABASE IF NOT EXISTS gobid;" || true

echo "applying migrations..."
cd /app/internal/store/pgstore/migrations
tern migrate

echo "starting server..."
cd /app

exec ./main
