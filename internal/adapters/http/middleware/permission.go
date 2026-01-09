package middleware

import (
	"encoding/json"
	"net/http"

	"horizonx/internal/domain"
)

func Permission(roleSvc domain.RoleService, perm domain.PermissionConst) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := roleSvc.HasPermission(r.Context(), perm); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"message": domain.ErrYouDontHavePermission.Error(),
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
