# KrakenD JWT Revoke Plugin

This repository hosts a plugin for [KrakenD](https://github.com/krakend/krakend-ce), the high-performance API Gateway. The plugin integrates a JWT Revoke Server seamlessly into KrakenD deployments, providing enhanced JWT token management capabilities.

## Usage

### Docker Image

To use the Docker Image, you can pull it from the following location:

```
ghcr.io/k-orolevsk-y/krakend-jwt-revoker
```

### Kubernetes Deployment

To integrate this Docker Image into a Kubernetes deployment, you can include it in your Kubernetes manifests.

### Initialization in krakend.json

To initialize the plugin in your KrakenD configuration file (krakend.json), follow these steps:

1. Add the following section to define the plugin directory:

    ```json
    {
      "plugin": {
        "pattern": ".so",
        "folder": "./plugins/"
      }
    }   
    ```

2. Specify the plugin in the `extra_config` section as follows:

    ```json
    {
      "extra_config": {
        "plugin/http-server": {
          "name": ["krakend-jwt-revoker"]
        }
      }
    }
    ```

### Configuration Options

The plugin supports the following configuration options:

- **addr**: The address of the HTTP revoke server (e.g., ":9000").
- **response**: A map of key-value pairs representing the response to be sent upon token revocation.
- **key**: The key where the token is expected to be found.
- **keyType**: The location where the token is expected to be found, with possible values of "header" or "cookie".

Example configuration:

```json
{
  "krakend-jwt-revoker": {
    "addr": ":9000",
    "response": {},
    "key": "token",
    "keyType": "header"
  }
}
```

Ensure that you adjust these configuration options according to your specific requirements and environment setup.

## Revoke Server Endpoint

The JWT Revoke Server exposes a single POST endpoint for token revocation:

```
POST http://{addr}/
```

### Request Body

The request body should be a JSON array containing objects with key-value pairs specifying the criteria for token revocation. Each object represents a criterion to be applied to the tokens. The following structure is expected:

```json
[
    {
        "key": "uuid",
        "value": "..."
    },
    {
        "key": "...",
        "value": "..."
    }
]
```

- **key**: The key specifying the criterion (e.g., "uuid" for the unique identifier).
- **value**: The corresponding value against which tokens will be evaluated. All values are converted to strings for comparison.


Tokens will be revoked based on the specified criteria. It's essential to note that each token must contain the specified fields for revocation and the iat field, which will be used for token evaluation by time.

## Contribution

Contributions are welcome! Feel free to submit issues or pull requests if you encounter any problems or have suggestions for improvements.

---

*This project is maintained by [k-orolevsk-y](https://github.com/k-orolevsk-y).*