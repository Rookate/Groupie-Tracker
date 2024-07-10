# Groupie Trackers API Web Application

## Overview

This project is a web application built with Go (Golang) that interacts with the Groupie Trackers API. The application fetches data about musical artists, including their details, locations, concert dates, and relations, and displays it to users. The project includes features such as viewing all artists, detailed views of each artist, and a search functionality.

## Features

- **Home Page**: Displays a list of all artists fetched from the Groupie Trackers API.
- **Detail Page**: Provides detailed information about a selected artist, including locations, concert dates, and relations.
- **Search**: Allows users to search for artists based on their name, members, first album, creation date, locations, and concert dates.
- **Error Handling**: Displays custom error pages for various error scenarios.
- **Static File Handling**: Serves static files (e.g., CSS, JavaScript) for the web application.

## Installation

1. **Clone the repository**:

   ```sh
   git clone https://zone01normandie.org/git/ejean/groupie-tracker.git
   cd groupie-trackers
   ```

2. **Install dependencies**:
   Make sure you have Go installed on your machine. If not, download and install it from [golang.org](https://golang.org/dl/).

3. **Run the application**:

   ```sh
   go run main.go
   ```

4. **Access the application**:
   Open your web browser and navigate to `http://localhost:8080`.

## Project Structure

- `main.go`: The main entry point of the application. It sets up the server and routes.
- `templates/`: Contains HTML templates for the home, detail, search, and error pages.
- `static/`: Contains static files like CSS and JavaScript.
- `Artist`: Struct to represent an artist fetched from the API.
- `Locations`: Struct to represent an artist's locations.
- `Dates`: Struct to represent an artist's concert dates.
- `Relation`: Struct to represent an artist's relations.
- `ArtistInfo`: Struct to aggregate detailed information about an artist.

## Endpoints

- **Home Page**: `/`
  - Displays a list of all artists.
- **Detail Page**: `/detail/{id}`
  - Displays detailed information about a specific artist by ID.
- **Search**: `/search?query={search_term}`
  - Allows users to search for artists.

## Error Handling

The application handles errors gracefully by displaying custom error pages for scenarios such as:
- Page not found (404)
- Internal server error (500)
- Issues with fetching data from the Groupie Trackers API

## Environment Variables

- `ApiURL`: The base URL of the Groupie Trackers API. Set to `https://groupietrackers.herokuapp.com/api`.

## Example

Here is an example of how to start the server:

```sh
$ go run main.go
Server starting on http://localhost:8080
```

Visit `http://localhost:8080` in your web browser to view the home page of the application.

---

By following this README, you should be able to set up and run the Groupie Trackers API Web Application on your local machine. Feel free to explore and enhance the application as needed.