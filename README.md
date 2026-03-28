# Portfolio Backend

A modern, scalable backend service for managing my portfolio website.

---

## Prerequisites

Ensure the following tools are installed:

- **Docker** (for containers and dependencies)
- **Make** (to execute Makefile commands)
- **Git** (for source control and pre-commit hooks)
- **Go** (v1.21+ recommended if running outside Docker)

---

## Setup

### 1. Configure Environment Variables

Copy the example environment file:

```bash
cp .default.env .local.env
```

Edit `.local.env` with your credentials, secrets, and any custom database/JWT/etc. settings.

### 2. (Optional) Install Pre-commit Hooks

Install code quality hooks:

```bash
pre-commit install
```

This runs code checks automatically before each commit.

### 3. Initialize Go Dependencies

```bash
make init
```

Installs Go modules and development dependencies.

---

## Running the Application

### 1. Start Backend Dependencies (Postgres)

```bash
make dockerup
```

This launches backend dependencies in containers.

### 2. Run the Backend Locally

This also applies database migrations automatically:

```bash
make run
```

(Default port: `8000`)

---

## Stopping & Cleaning Up

- **Gracefully stop backend:** Press <kbd>Ctrl</kbd>+<kbd>C</kbd> in your app terminal.
- **Tear down Docker dependencies:**

  ```bash
  make dockerdown
  ```

  This shuts down and removes all running containers.

---

## Running Tests

Run all tests:

```bash
make test
```

Run specific tests by package or file:

```bash
go test -v -run TestName ./path/to/package_or_file.go
```

---

## Creating and Managing Migrations

After updating your Go models, generate a new migration script:

```bash
make migrate-create "description"
```

#### Migration Naming Convention

Replace `"description"` with a concise, snake_case summary:

- **Single Table**: `<table_name>_<action>_<target>`
- **Multiple Tables**: `feature_<general_summary>`

Where:

- `<table_name>`: Plural name of the primary table (put first for grouping)
- `<action>`: SQL verb (create, add, drop, alter, rename)
- `<target>`: Column/index/constraint being changed; for new tables, use `table`
- `feature_`: Prefix if touching 3+ tables
- `<general_summary>`: Concise summary for multi-table migrations

**Examples:**

_Single Table:_

- `users_create_table` &nbsp;&nbsp;// becomes: 001_users_create_table
- `users_add_email` &nbsp;&nbsp;// becomes: 002_users_add_email
- `users_alter_email_length`
- `users_drop_email`
- `users_add_index_email`

_Multi-Table:_

- `feature_user_onboarding_setup`
- `feature_stripe_billing_integration`

This scheme makes every migration script clear in intent and scope.

---

## Resetting the Database (Development Only)

Wipe and re-migrate all tables:

```bash
make db-fresh
```

---

## Contribution Policy

This codebase powers my personal website and is provided here solely as a portfolio example.  
I am not accepting pull requests, issues, or external contributions at this time.

---

Copyright (c) 2026 Fabian Railey Victuelles. All Rights Reserved.
