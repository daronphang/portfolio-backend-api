## Config

Requires env.yaml, gmail.credentials.json, gmail.token.json.

### Gmail

Requires refresh token to be present. Generated on first request to Google server.

https://developers.google.com/gmail/api/quickstart/go

### Refresh Token

Doesn't expire unless certain condtions are met:

- User has revoked access to the application.
- Refresh token has not been used for six months.
- User changed password and the refresh token contained Gmail scopes.
