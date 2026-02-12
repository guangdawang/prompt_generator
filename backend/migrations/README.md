# Database Migrations

This directory contains SQL migration files for the database schema and initial data.

## Migration Files

- `001_initial_schema.sql` - Creates the initial database schema (tables, indexes)
- `002_seed_data.sql` - Inserts sample data for development

## How Migrations Work

1. Migration files are executed in alphabetical order
2. Files ending with `.sql` are automatically detected and executed
3. Each migration is run once when the application starts
4. Use `IF NOT EXISTS` and `ON CONFLICT` clauses for idempotent operations

## Adding New Migrations

1. Create a new SQL file with the next sequential number (e.g., `003_new_feature.sql`)
2. Write your SQL statements
3. Test the migration locally
4. Commit the migration file

## Migration Best Practices

- Use descriptive file names
- Make migrations idempotent (can be run multiple times safely)
- Test migrations on a copy of production data before applying
- Include rollback scripts for complex changes (when needed)
