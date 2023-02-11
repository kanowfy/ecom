<!-- PROJECT LOGO -->
<br />
<div align="center">

  <h3 align="center">ECOM</h3>

  <p align="center">
  A Demo Ecommerce Project built using P2P Microservice Architecture for learning purposes
  </p>
</div>

<!-- ABOUT THE PROJECT -->

## About The Project

<img src="https://raw.githubusercontent.com/kanowfy/ecom/dev/img/1.png">

8 Go services communicate through <b>gRPC</b>:
<br/>

- <b>Catalog Service</b>: Handle management of products, uses Postgresql to persist data
- <b>Cart Service</b>: Handle customers' shopping carts, uses Redis to persist data
- <b>Payment Service</b>: Mock handling payment
- <b>Shipping Service</b>: Get shipping cost and mock shipping order
- <b>Email Service</b>: Send confirmation email to customer, uses [Mailtrap](https://mailtrap.io) to test inbox
- <b>Order Service</b>: Aggregate data to produce order result
- <b>Frontend Service</b>: Display the UI

### Prerequisites

- Docker with Docker-Compose installed

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/kanowfy/ecom.git
   ```
2. Run docker-compose
   ```sh
   docker compose up
   ```
3. Visit `localhost:4000`

<!-- CONTACT -->

## Contact

Your Name - [@kanowfy](https://twitter.com/kanowfy) - nguyendung2002hl@gmail.com

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

- [Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo)
- [Docker Compose Documentation](https://docs.docker.com/engine/reference/commandline/compose_up/)
