package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
)

func checkSessionID(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		c, err := r.Cookie(cookieSID)
		if errors.Is(err, http.ErrNoCookie) {
			u, _ := uuid.NewV4()
			sessionID = u.String()
			http.SetCookie(w, &http.Cookie{
				Name:   cookieSID,
				Value:  sessionID,
				MaxAge: cookieMaxAge,
			})
		} else if err != nil {
			return
		} else {
			sessionID = c.Value
		}

		ctx := context.WithValue(r.Context(), ctxSIDKey{}, sessionID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
