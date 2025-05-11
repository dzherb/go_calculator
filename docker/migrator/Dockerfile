FROM migrate/migrate

COPY ./calculator/internal/storage/migrations ./migrations

ENTRYPOINT migrate -path /migrations -database $DATABASE_URL up