# Gator CLI

Gator is a powerful command-line interface tool designed to help you manage and browse RSS feeds right from your terminal. Stay up-to-date with your favorite content without ever leaving the command line!

## Prerequisites

To get Gator up and running on your system, you'll need two main components:

*   **Go**: Gator is built with Go. You'll need a Go development environment installed on your machine. You can find detailed installation instructions on the official [Go website](https://golang.org/doc/install).
*   **PostgreSQL**: Gator uses a PostgreSQL database to store all your feed and post data. You'll need to have a PostgreSQL server running and know its connection details (username, password, host, port, database name). If you don't have PostgreSQL installed, you can find installation guides on the official [PostgreSQL website](https://www.postgresql.org/download/).

## Installation

Once you have Go installed, installing Gator is quite straightforward:

1.  **Clone the Gator repository:**
    Open your terminal and clone the project to your local machine:
    ```bash
    git clone https://github.com/YOUR_GITHUB_USERNAME/YOUR_REPO_NAME.git
    cd YOUR_REPO_NAME
    ```
    (Remember to replace `YOUR_GITHUB_USERNAME` and `YOUR_REPO_NAME` with the actual path to your repository!)

2.  **Install the Gator CLI tool:**
    Navigate into the cloned directory and run the `go install` command:
    ```bash
    go install .
    ```
    This command compiles the Gator application and places the `gator` executable in your Go binary path (typically `~/go/bin`), making it available from any directory in your terminal.

## Configuration and Running

Gator requires a small configuration file to connect to your database.

1.  **Create a `.env` file:** In the root directory of the `GatorRSS` project (where you cloned it), create a new file named `.env`.

2.  **Add your database connection string:** Inside the `.env` file, add a single line specifying your PostgreSQL database URL. Replace the placeholders with your actual database credentials:
    ```
    DATABASE_URL=postgres://user:password@host:port/database_name?sslmode=disable
    ```
    For example: `DATABASE_URL=postgres://gator_user:mysecretpassword@localhost:5432/gator_db?sslmode=disable`

3.  **Run database migrations:** Before using Gator for the first time, you need to set up the database schema. Run the `migrate` command:
    ```bash
    gator migrate
    ```

4.  **Start using Gator:** Now you're ready to use the Gator CLI!

## Commands

Here's a list of the main commands you can use with Gator:

### User Management

*   **`gator register`**: Create a new user account to get started.
*   **`gator login`**: Log in to your Gator account to access personalized features.
*   **`gator reset`**: Resets your user password.
*   **`gator users`**: Lists all registered users (useful for administration or testing).

### Feed Management and Browsing

*   **`gator addfeed <URL>`**: Add a new RSS feed to your collection. (Requires login)
*   **`gator feeds`**: List all the feeds that have been added to the system.
*   **`gator follow <FeedID>`**: Start following a specific feed by its ID to receive its posts. (Requires login)
*   **`gator following`**: See a list of all the feeds you are currently following. (Requires login)
*   **`gator unfollow <FeedID>`**: Stop following a feed by its ID. (Requires login)
*   **`gator agg`**: This command is likely used for aggregation or fetching, though its exact behavior might vary based on your implementation details.
*   **`gator browse [limit]`**: View the latest posts from all the feeds you follow. You can optionally specify a `limit` to control how many posts are displayed (e.g., `gator browse 5` to show the 5 most recent posts). (Requires login)
