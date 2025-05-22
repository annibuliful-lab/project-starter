package authentication

import (
	error_utils "backend/src/error"
	"backend/src/jwt"
	"context"
	"encoding/json"
	"net/http"
)

func GetAuthorizationHeaders(header http.Header) (AuthorizationContext, error) {
	authorization := header.Get("authorization")
	projectId := header.Get("x-project-id")

	authContext := AuthorizationContext{}

	if authorization == "" && projectId != "" {
		return authContext, nil
	}

	if projectId != "" {
		authContext.ProjectId = projectId
	}

	authContext.Token = GetAuthToken(authorization)

	if authContext.Token == "" {
		return authContext, nil
	}

	payload, err := jwt.VerifyToken(authContext.Token)

	if err != nil {
		return authContext, err
	}

	authContext.AccountId = payload.AccountId.String()

	return authContext, nil
}

func GraphqlContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers, err := GetAuthorizationHeaders(r.Header)

		if headers.Token == "" {
			next.ServeHTTP(w, r)
			return
		}

		payload, err := jwt.VerifyToken(headers.Token)
		if err != nil {
			if err.Error() == error_utils.TokenExpire.Error() {
				// Token expired error
				errorResponse := map[string]string{"error": error_utils.TokenExpire.Error()}
				writeJSONResponse(w, errorResponse, http.StatusUnauthorized)
				return
			}
			// Other token verification errors
			errorResponse := map[string]string{"error": error_utils.TokenVerificationFailed.Error()}
			writeJSONResponse(w, errorResponse, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", headers.Token)

		if headers.ProjectId != "" {
			ctx = context.WithValue(ctx, "projectId", headers.ProjectId)
		}

		ctx = context.WithValue(ctx, "accountId", payload.AccountId.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
