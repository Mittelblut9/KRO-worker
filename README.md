# KommtRumathraOnline.de Worker

This is the worker for kommtrumathraonline.de
The worker runs each hour and downloads new vods if found.
After the download it gets transcribed with openai whisper.
The result together with an propmt goes to ChatGPT which then decides if Rumathra has intend of coming online.
ChatGPT returns a json, the json and the vod info from twitch will be written to the database.

## Usage

1. Copy the .env.example to .env
2. Insert missing values
3. Start the project:

```bash
docker-compose up
```

## How does it work

![image](docs/flow.png)

## Original

[KommtKevinOnline Worker](https://github.com/KommtKevinOnline/worker)
