# Tap It - Text Phishing Framework

Tap It is a SMS phishing framework that allows handling of large SMS phishing campaigns. It automatically handles text template and allows basic handling of web template should phishing URLs be sent through the text.

### Prerequisites

Tap It is built on Go, and have all required libraries built-in within the binary.

To attempt to recompile the application, the following is required:

* Go
* Angular CLI
* GORM
* Gorilla Mux
* Teabag XLSX Library

## Deployment

The entire application is designed to be Dockerised. To build the Docker environment with the pre-compiled binaries, simple run the following:
```
sudo docker-compose up
```
You may access the management platform at http://127.0.0.1:8000/
Victim access is at http://127.0.0.1:8001/

Default registration password for management platform is "Super-Secret-Code"

## Authors

* **Samuel Pua** - *Initial work* - [GitHub](https://github.com/telboon)

## License

This project is licensed under the BSD License - see the [LICENSE](LICENSE) file for details


