# File Upload and Download Application

This is a simple web application that allows users to upload files and download them. It provides a basic frontend interface for file uploads and displays a list of uploaded files. Users can click on the filenames to download the corresponding files.

## Prerequisites

- Go programming language (version 1.16+)
- Git
- Docker (Optional)

## Getting Started

1. Clone the repository:
git clone <repository_url>

2. Navigate to the project directory:
```cd go-upload```
3. Build the application:
```go build```

## Getting Started Docker Version
1. Build Docker image
```docker build -t go-upload .```
2. Run Docker container
````docker run -p 8000:8000 go-upload```
4. Run the application:

5. Open your web browser and visit `http://localhost:8000` to access the application.

## Usage

- To upload a file, click on the "Choose File" button, select a file from your local system, and click the "Upload" button. The progress bar will show the upload progress.

- The uploaded files will be listed below the upload form. Click on the filename to download the corresponding file.

## Configuration

- The application uses the `uploads` directory to store uploaded files. Make sure this directory exists and is writable by the application.

## License

This project is licensed under the [MIT License](LICENSE).



