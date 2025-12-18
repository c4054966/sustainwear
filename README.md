# Sustainwear

A full-stack sustainability-focused application with a Go backend and Next.js frontend.

## Backend Setup (GO)

### Prerequisites

Ensure you have Go installed on your system.

### Installation

Navigate to the backend directory:

cd backend


Download the required dependencies:

go mod download

or

go mod tidy


### Running the Backend

You have two options to start the backend server:

#### Option 1: Run directly in terminal

go run main.go


#### Option 2: Build and run executable

go build -o Sustainwear.exe


Then execute the generated file:

./Sustainwear.exe


### Configuration

The backend uses `config.toml` for all configuration settings. A default configuration file is included in the repository.

If the config file is missing or deleted, simply restart the program to automatically generate a new `config.toml` with default values.

**Important:** Configuration changes require a backend restart to take effect, as the config is loaded on startup.

By default, the application is configured to use SQLite, but you can easily switch to other database drivers by modifying the values in `config.toml`.

## Frontend Setup (NextJS)

### Prerequisites

Ensure you have Node.js and npm installed on your system.

### Installation

Navigate to the frontend directory:

cd frontend

Install the required dependencies:

npm install
npm install leaflet react-leaflet
npm i --save-dev @types/leaflet


### Running the Frontend

Start the development server:

npm run dev


### Configuration

The frontend configuration is managed in `next.config.ts`, which specifies the backend destination URL.

You can modify this file if you're using custom backend configuration values different from the defaults.