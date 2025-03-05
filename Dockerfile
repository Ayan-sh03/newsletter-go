# Use the official lightweight Postgres Alpine image
FROM postgres:16-alpine

# Set environment variables
ENV POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=postgres \
    POSTGRES_DB=newsletter

# Expose the default PostgreSQL port
EXPOSE 5432

# Use a volume for persistent storage (defined in docker-compose or runtime)
VOLUME ["/var/lib/postgresql/data"]

# Set the default command to run Postgres
CMD ["postgres"]
