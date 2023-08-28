package types

type ZitadelError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type IntrospectionResult struct {
	Active    bool     `json:"active"`
	Aud       []string `json:"aud"`
	ClientId  string   `json:"client_id"`
	Exp       int      `json:"exp"`
	Iat       int      `json:"iat"`
	Iss       string   `json:"iss"`
	Jti       string   `json:"jti"`
	Nbf       int      `json:"nbf"`
	Scope     string   `json:"scope"`
	TokenType string   `json:"token_type"`
	Username  string   `json:"username"`
	Sub       string   `json:"sub"`
	Roles     any      `json:"urn:zitadel:iam:org:project:roles"`
}
