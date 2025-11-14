# Swapdoodle replacement server

A HPP server for the _Swapdoodle_ 3DS app.

## Configuration

All configuration options are handled via environment variables

`.env` files are supported

| Name                            | Description                                    | Required                            |
| ------------------------------- | ---------------------------------------------- | ----------------------------------- |
| `PN_SD_POSTGRES_URI`            | Fully qualified URI to your Postgres server    | Yes                                 |
| `PN_SD_HPP_SERVER_PORT`         | Port for the HPP server                        | Yes                                 |
| `PN_SD_GRPC_SERVER_PORT`        | Port for the GRPC server                       | Yes                                 |
| `PN_SD_CONFIG_S3_ENDPOINT`      | S3 server endpoint                             | Yes                                 |
| `PN_SD_CONFIG_S3_ACCESS_KEY`    | S3 access key ID                               | Yes                                 |
| `PN_SD_CONFIG_S3_ACCESS_SECRET` | S3 secret                                      | Yes                                 |
| `PN_SD_CONFIG_S3_BUCKET`        | S3 bucket                                      | Yes                                 |
| `PN_SD_CONFIG_GRPC_API_KEY`     | API Key                                        | No (in which case it is open)       |
| `PN_SD_ACCOUNT_GRPC_HOST`       | Host name for your account server gRPC service | Yes                                 |
| `PN_SD_ACCOUNT_GRPC_PORT`       | Port for your account server gRPC service      | Yes                                 |
| `PN_SD_ACCOUNT_GRPC_API_KEY`    | API key for your account server gRPC service   | No (Assumed to be an open gRPC API) |
