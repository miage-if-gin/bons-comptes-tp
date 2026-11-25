package http

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/miage-if-gin/bons-comptes-tp/internal/service"
)

// Ce test ne touche ni la base ni le serveur : il appelle directement la
// fonction de traduction et lit le code écrit dans la réponse.
//
// Il ne couvre que le TODO 1.3. Les TODO 1.1 et 1.2 se vérifient dans le
// navigateur : la page de clôture doit s'afficher.
func TestRepondreErreurCloture(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cas := []struct {
		nom     string
		err     error
		attendu int
	}{
		{"un porte-monnaie qui n'existe pas", service.ErrPorteMonnaieIntrouvable, 404},
		{"un virement qui n'existe pas", service.ErrVirementIntrouvable, 404},
		{"on n'est pas membre", service.ErrAccesRefuse, 403},
		{"seul le créateur peut clôturer", service.ErrClotureReservee, 403},
		{"seul le débiteur peut rembourser", service.ErrRemboursementInterdit, 403},
		{"déjà clos", service.ErrDejaClos, 409},
		{"pas encore clos", service.ErrPasEncoreClos, 409},
		{"déjà remboursé", service.ErrDejaRembourse, 409},
		{"une panne imprévue", errors.New("base injoignable"), 500},
	}

	for _, c := range cas {
		t.Run(c.nom, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest("GET", "/", nil)

			repondreErreurCloture(ctx, c.err)

			if rec.Code != c.attendu {
				t.Fatalf("« %s » → code %d, attendu %d\n\nerreur : %v",
					c.nom, rec.Code, c.attendu, c.err)
			}
		})
	}
}
