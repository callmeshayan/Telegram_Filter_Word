# Telegram Filter Bot

This bot filters messages containing a specified word and logs them into a PostgreSQL database.

## Prerequisites

1. Go programming language installed: https://golang.org/doc/install
2. PostgreSQL database set up and running
3. Telegram Bot Token: Obtain from BotFather on Telegram
4. Go modules installed: Run `go mod init` in your project directory

## Setup

1. Clone this repository.
2. Set up a PostgreSQL database and note down the database URL and password.
3. Set the environment variables:
    - `BOT_TOKEN`: Your Telegram Bot Token
    - `POSTGRES_PASSWORD`: Your PostgreSQL password
4. Build the bot: Run `go build` in the project directory.
5. Run the compiled binary.

## Commands

### /start
- Description: Start the bot.
- Usage: `/start`

### /filter
- Description: Filter messages containing a specified word.
- Usage: `/filter 'word'`
- Example: `/filter example`
- Note: Replace `'word'` with the actual word you want to filter.

### /help
- Description: Display help message.
- Usage: `/help`

### /stop
- Description: Stop the bot.
- Usage: `/stop`

## Functionality

- **Filtering Messages**: Use the `/filter` command followed by a word to filter messages containing that word.
- **Help Message**: Use the `/help` command to display available commands and usage.
- **Stopping the Bot**: Use the `/stop` command to stop the bot.
- **Logging**: Filtered messages are logged into the `filtered_messages` table while non-filtered messages are logged into the `non_filtered_messages` table in the PostgreSQL database.

## Implementation Details

- **Logging Messages**: Messages filtered by the bot are inserted into the `filtered_messages` table, and non-filtered messages are inserted into the `non_filtered_messages` table in the PostgreSQL database.
- **Database Schema**:
  - `filtered_messages` table: 
    - Columns: `message_id` (INT), `chat_id` (INT), `user_id` (INT), `message_text` (TEXT), `timestamp` (TIMESTAMP)
  - `non_filtered_messages` table:
    - Columns: `message_id` (INT), `chat_id` (INT), `user_id` (INT), `message_text` (TEXT), `timestamp` (TIMESTAMP)
- **Database Connection**: The bot establishes a connection to the PostgreSQL database using the provided database URL and password. 
- **Error Handling**: Errors related to database operations are logged, and appropriate error messages are sent to users.
