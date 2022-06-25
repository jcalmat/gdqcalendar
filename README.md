# Games Done Quick Calendar

GDQ Calendar provides some tools to extract events from the Games Done Quick [schedule](https://www.gamesdonequick.com/schedule).

It allows you to generate an ICS file in order to import it by yourself, or to connect your own calendar and directly import the events inside it.

## 1. Generate an ICS file

In order to generate an ICS file, you need to run the following command:

```golang
go run cmd/generator/main.go
```

The `gdq.ics` file will be generated in the current directory and can be imported into your calendar application.


## 2. Use Google API to directly fill your calendar with the next GDQ events

This project allows you to connect your own Google cloud project in order to directly import the current GDQ events into your calendar.

### Setup

1. Create a [google cloud project](https://developers.google.com/workspace/guides/create-project)
2. Create an [OAuth access credentials](https://developers.google.com/workspace/guides/create-credentials#oauth-client-id)
3. Put your credentials in a `credentials.json` file in the current directory
4. Create a `config.json file` and fill it with the following content:
    ```json
    {
      "calendarId": "your-calendar-id"
    }
    ```
    ([How to find my calendar ID](https://fullcalendar.io/docs/google-calendar))

### Usage

```golang
    go run cmd/api/*.go
```


**First launch**

The first time you'll run this project, you'll need to grant access to your calendar.

The program will prompt you a link asking for your authorization.

Accept, you will be given a code that you have to paste back to this program.

This will create a `token.json` file that will be used for authentication and authorization.


**Subsequent launches**

This program stores the events it's creating. In order to avoid any event duplication, each time you'll rerun the program, it will remove all the stored events and regenereate them.