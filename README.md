# Games Done Quick Calendar

GDQ Calendar provides tools to extract events from the Games Done Quick [schedule](https://www.gamesdonequick.com/schedule).

It enables you to generate an ICS file for manual import or to connect your calendar directly and import events automatically.

## 1. Generate an ICS file

To generate an ICS file, run the following command:

```bash
go run cmd/generator/main.go
```

The `gdq.ics` file will be created in the current directory for import into your calendar application.

## 2. Use Google API to Automatically Populate Your Calendar with Upcoming GDQ Events

This project allows you to connect your Google Cloud project to import current GDQ events directly into your calendar.

### Setup

1. Create a [Google Cloud project](https://developers.google.com/workspace/guides/create-project).
2. Generate [OAuth access credentials](https://developers.google.com/workspace/guides/create-credentials#oauth-client-id).
3. Save your credentials in a `credentials.json` file in the current directory.
4. Create a `config.json` file and fill it with the following content:

```json
{
  "calendarId": "your-calendar-id"
}
```

([How to find my calendar ID](https://fullcalendar.io/docs/google-calendar))

### Usage

```bash
go run cmd/api/*.go
```

**First Launch**

The first time you run this project, you'll need to grant access to your calendar. The program will provide a link for your authorization. Accept, and you will receive a code to paste back into the program. This creates a `token.json` file for authentication and authorization.

**Subsequent Launches**

Upon every launch, the program refreshes upcoming events by removing them through the Google API and recreating them. This ensures that the events are up-to-date.
