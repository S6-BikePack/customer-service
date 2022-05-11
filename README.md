<div align="center">
  <a href="https://github.com/S6-BikePack">
    <img src="assets/logo.png" alt="logo" width="200" height="auto" />
  </a>
  <h1>BikePack - Customer-Service</h1>

  <p>
    Part of the S6 BikePack project.
  </p>


<!-- Badges -->
<p>

</p>

<h4>
    <a href="https://github.com/S6-BikePack">Home</a>
  <span> · </span>
    <a href="https://github.com/S6-BikePack/customer-service#-about-the-project">Documentation</a>
  </h4>
</div>

<br />

<!-- Table of Contents -->
# 📓 Table of Contents

- [About the Project](#-about-the-project)
    * [Tech Stack](#%EF%B8%8F-tech-stack)
    * [Environment Variables](#-environment-variables)
- [Getting Started](%EF%B8%8F-getting-started)
    * [Prerequisites](%EF%B8%8F-prerequisites)
    * [Running Tests](#-running-tests)
    * [Run Locally](#-run-locally)
    * [Deployment](#-deployment)
- [Usage](#-usage)



<!-- About the Project -->
## ⭐ About the Project

The Customer-Service is the service for the BikePack project that handles all customers in the system.


<!-- TechStack -->
### 🛰️ Tech Stack
#### Language
  <ul>
    <li><a href="https://go.dev/">GoLang</a></li>
</ul>

#### Dependencies
  <ul>
    <li><a href="https://github.com/gin-gonic/gin">Gin</a><span> - Web framework</span></li>
    <li><a href="https://github.com/swaggo/swag">Swag</a><span> - Swagger documentation</span></li>
    <li><a href="https://gorm.io/index.html">GORM</a><span> - ORM library</span></li>
  </ul>

<!-- Env Variables -->
### 🔑 Environment Variables

This service has the following environment variables that can be set:

`PORT` - Port the service runs on
`DATABASE` - Databae connection url
`RABBITMQ` - RabbitMQ connection url

<!-- Getting Started -->
## 	🛠️ Getting Started

<!-- Prerequisites -->
### ‼️ Prerequisites

Building the project requires Go 1.18.

The easiest way to setup the project is to use the Docker-Compose file from the infrastructure repository.

<!-- Running Tests -->
### 🧪 Running Tests

-

<!-- Run Locally -->
### 🏃 Run Locally

Clone the project

```bash
  git clone https://github.com/S6-BikePack/customer-service
```

Go to the project directory

```bash
  cd customer-service
```

Run the project (Rest)

```bash
  go run cmd/rest/main.go
```


<!-- Deployment -->
### 🚀 Deployment

To build this project run (Rest)

```bash
  go build cmd/rest/main.go
```


<!-- Usage -->
## 👀 Usage

### REST
Once the service is running you can find its swagger documentation with all the endpoints at `/swagger`