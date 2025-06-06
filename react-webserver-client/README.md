# React Webserver Client

This project is a React application that interacts with a Go webserver API to display information about people. It includes components for listing people, displaying details about a specific person, and greeting users.

## Project Structure

```
react-webserver-client
├── public
│   ├── index.html        # Main HTML file for the React application
│   └── favicon.ico       # Favicon for the application
├── src
│   ├── components        # Contains React components
│   │   ├── PeopleList.jsx  # Component to display a list of people
│   │   ├── PersonDetail.jsx # Component to display details of a specific person
│   │   └── Greeting.jsx     # Component to display a greeting message
│   ├── services          # Contains API service functions
│   │   └── api.js       # Functions for making API calls to the Go webserver
│   ├── App.jsx           # Main application component
│   ├── index.js          # Entry point for the React application
│   └── styles            # Contains CSS styles
│       └── App.css      # Styles for the application
├── package.json          # Configuration file for npm
├── .gitignore            # Specifies files to ignore by Git
└── README.md             # Documentation for the project
```

## Getting Started

To get started with the project, follow these steps:

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd react-webserver-client
   ```

2. **Install dependencies**:
   ```
   npm install
   ```

3. **Run the application**:
   ```
   npm start
   ```

The application will be available at `http://localhost:3000`.

## API Endpoints

The application interacts with the following API endpoints provided by the Go webserver:

- `GET /people`: Fetches a list of people.
- `GET /people/:id`: Fetches details of a specific person by ID.
- `GET /`: Returns a greeting message.
- `GET /greet/:name`: Returns a personalized greeting message.

## Components Overview

- **PeopleList**: Fetches and displays a list of people from the API.
- **PersonDetail**: Fetches and displays detailed information about a specific person.
- **Greeting**: Displays a greeting message.

## License

This project is licensed under the MIT License.