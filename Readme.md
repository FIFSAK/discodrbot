# Discord Voice Room Auto-Recording Bot with S3 Storage

This bot is designed to automatically record new voice rooms in Discord when at least one user enters the room, and store the recordings in an S3 bucket. The recording stops when the time limit is reached or the voice channel becomes empty. After recording, the bot sends a link to the recorded file in the voice channel's chat. The bot is written in Go and uses the `discordgo` package to interact with Discord and the `aws-sdk-go` package to interact with AWS S3.

## Documentation

- **bot.go** - Main file with bot logic.
- **managefile.go** - Logic for saving PCM to file and uploading to S3.
- **record.go** - Logic for recording voice channels.
- **s3.go** - Logic for uploading files to S3.

## Configuration

Create a `.env` file with the following environment variables:

```
DISCORD_BOT_TOKEN=<Your Discord Bot Token>
S3_API_KEY=<Your AWS S3 API Key>
S3_SECRET_KEY=<Your AWS S3 Secret Key>
S3_REGION=<Your AWS S3 Region>
S3_BUCKET_NAME=<Your AWS S3 Bucket Name>
```

## How to Run

1. **Clone the repository:**

   ```sh
   git clone https://github.com/FIFSAK/discodrbot.git
   cd discodrbot
   ```

2. **Initialize Go modules and install dependencies:**

   ```sh
   go mod tidy
   ```

   This command will download all the necessary dependencies specified in the `go.mod` file.

```sh
go run cmd/main.go
```



