# Login Microservice with Keycloak Integration

This project demonstrates a login microservice built in Go, integrated with Keycloak for authentication. It supports two OAuth2 flows:
- `grant_type=password` (Resource Owner Password Credentials Grant)
- `grant_type=authorization_code` with PKCE (Proof Key for Code Exchange)

## Comparison: `grant_type=password` vs `grant_type=authorization_code`

| Aspect                        | `grant_type=password`                          | `grant_type=authorization_code` with PKCE |
|-------------------------------|-----------------------------------------------|------------------------------------------|
| **Security**                  | Less secure. Sends user credentials directly to the client. | More secure. Credentials are handled by Keycloak directly. |
| **Client Type**               | Suitable for confidential clients (e.g., backend servers). | Designed for public clients (e.g., mobile apps, SPAs). |
| **Client Secret**             | Requires a client secret.                     | Does not require a client secret.        |
| **Man-in-the-Middle Protection** | Vulnerable if HTTPS is not used.              | Protects against interception with `code_verifier` and `code_challenge`. |
| **OAuth2 Recommendation**     | Deprecated in modern OAuth2 implementations.  | Recommended by OAuth2.1 and widely adopted. |
| **Multi-Factor Authentication (MFA)** | Not supported natively.                     | Fully supported.                         |

## How to Use This Project

### Prerequisites
- Docker and Docker Compose installed.
- Keycloak running (configured via Docker Compose in this project).

### Setup
1. Start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. Access Keycloak at `http://localhost:8080`:
   - Username: `admin`
   - Password: `admin`

### Login as Admin to Get Token
You need to log in as the Keycloak admin to obtain an admin token. Use the following `curl` command:
```bash
curl -X POST "http://localhost:8080/realms/master/protocol/openid-connect/token" \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "username=admin" \
-d "password=admin" \
-d "grant_type=password" \
-d "client_id=admin-cli"
```
The response will include an `access_token`. Use this token in the `Authorization` header for subsequent requests.

### Configure a client in Keycloak:
   - Use the `/create-client` endpoint of this API to create and configure the client.
   - Example request:
     ```bash
     curl -X POST "http://localhost:8005/create-client" \
     -H "Authorization: Bearer <admin-token>" \
     -H "Content-Type: application/json" \
     -d '{
       "clientId": "client-test",
       "publicClient": true,
       "redirectUris": ["http://localhost:8005/redirect"],
       "protocol": "openid-connect"
     }'
     ```
   - Replace `<admin-token>` with a valid Keycloak admin token.

3. You can access Keycloak at the URL `http://localhost:8080`.

4. Create a user in Keycloak:

### Testing `grant_type=password`
1. Use the following `curl` command to test the `grant_type=password` flow:
   ```bash
   curl -X POST "http://localhost:8080/realms/master/protocol/openid-connect/token" \
   -H "Content-Type: application/x-www-form-urlencoded" \
   -d "grant_type=password" \
   -d "client_id=client-test" \
   -d "username=<username>" \
   -d "password=<password>"
   ```
2. The response will include an `access_token` and `refresh_token`.

### Testing `grant_type=authorization_code` with PKCE
1. Generate the `code_verifier` and `code_challenge`:
   ```bash
   curl -X POST "http://localhost:8005/generate-pkce" \
   -H "Content-Type: application/json" \
   -d '{
     "redirect_uri": "http://localhost:8005/redirect",
     "client_id": "client-test"
   }'
   ```
2. Copy the `authorization_url` from the response and open it in a browser.
3. Log in with the user credentials.
4. After login, you will be redirected to `http://localhost:8005/redirect` with the `code` in the query string.
5. Exchange the `code` for tokens:
   ```bash
   curl -X POST "http://localhost:8080/realms/master/protocol/openid-connect/token" \
   -H "Content-Type: application/x-www-form-urlencoded" \
   -d "grant_type=authorization_code" \
   -d "client_id=client-test" \
   -d "code=<authorization_code>" \
   -d "redirect_uri=http://localhost:8005/redirect" \
   -d "code_verifier=<code_verifier>"
   ```
6. The response will include an `access_token` and `refresh_token`.

### Custom Theme
This project includes a custom Keycloak theme located in the `keycloak-theme` directory. The theme is automatically applied when the Keycloak container is started.

To modify the theme, edit the files in the `keycloak-theme` directory and restart the Keycloak container:
```bash
docker-compose up --force-recreate keycloak
```

### Notes
- Always use HTTPS in production to secure communication.
- The `grant_type=password` flow is not recommended for public clients due to security concerns.
- The `grant_type=authorization_code` with PKCE is the preferred flow for modern applications.

Enjoy exploring the project!