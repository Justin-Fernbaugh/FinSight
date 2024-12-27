# FinSight

FinSight is a service designed to integrate with the budgeting software YNAB (You Need A Budget) and provide financial summaries via Telegram messages. The service is currently intended to run as a job on a schedule, although future enhancements will support processing many budgets from various users.

## Features

- Integrates with YNAB to retrieve transaction data.
- Uses Google Cloud's Vertex AI to generate financial summaries.
- Sends summaries to users via Telegram.
- Designed to run as a scheduled job.

## Prerequisites

- Go 1.22 or later
- Docker
- YNAB account and API token
- Google Cloud project with Vertex AI enabled
- Telegram bot token

## Environment Variables

The following environment variables need to be set:

- `PROJECT_ID`: The GCP project ID
- `LOCATION`: The GCP location (default: `us-west1`)
- `MODEL_NAME`: The LLM model name (default: `gemini-1.5-flash-001`)
- `YNAB_CLIENT_ID`: The YNAB application client ID
- `YNAB_TOKEN`: The YNAB client secret
- `DATABASE_NAME`: The GCP Firestore database name
- `TG_BOT_TOKEN`: The Telegram bot token
- `TG_USER_ID`: The Telegram user ID to send the message

## Build and Run

### Using Makefile

To build the Go binary:

```sh
make build
```

## Using Docker

To build the Docker image:
```Dockerfile
docker build -t finsight:latest .
```

To run the Docker container:
```Dockerfile
docker run -p 8080:8080 \
  -e PROJECT_ID="your_project_id" \
  -e LOCATION="us-west1" \
  -e MODEL_NAME="gemini-1.5-flash-001" \
  -e YNAB_CLIENT_ID="your_ynab_client_id" \
  -e YNAB_TOKEN="your_ynab_token" \
  -e DATABASE_NAME="your_database_name" \
  -e TG_BOT_TOKEN="your_tg_bot_token" \
  -e TG_USER_ID="your_tg_user_id" \
  finsight:latest
```

## How It Works
1. Initialization: The service initializes by setting up the required flags and environment variables.
2. Bot Handler: A Telegram bot handler is created to manage communication with the user.
3. LLM Handler: A handler for Google Cloud's Vertex AI is created to generate financial summaries.
4. YNAB Integration: The service retrieves transaction data from YNAB for the specified number of days back.
5. Summary Generation: The transaction data is processed and summarized using Vertex AI.
6. Message Sending: The generated summary is sent to the user via Telegram.

## Future Enhancements
- Support for processing multiple budgets from various users.
- Enhanced error handling and logging.
- Improved scheduling and job management.