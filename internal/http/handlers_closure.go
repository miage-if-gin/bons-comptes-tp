package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miage-if-gin/bons-comptes-tp/internal/service"
)

// registerClosure déclare les routes du lot D donc les  soldes, clôture et
// remboursements. Toutes dans le groupe privé, derrière exigerCompte().
func registerClosure(prive *gin.RouterGroup, d Deps) {
	prive.GET("/wallets/:id/closure", func(c *gin.Context) {
		walletID, ok := lireID(c, "id")
		if !ok {
			return
		}
		afficherCloture(c, d, walletID, http.StatusOK, "")
	})

	prive.POST("/wallets/:id/closure", func(c *gin.Context) {
		walletID, ok := lireID(c, "id")
		if !ok {
			return
		}
		compte := CompteCourant(c)

		if err := d.Closure.Cloturer(c.Request.Context(), walletID, compte.ID); err != nil {
			if estRefusCloture(err) {
				// On réaffiche la page avec le message plutôt qu'une page d'erreur nue.
				status := http.StatusConflict
				if errors.Is(err, service.ErrClotureReservee) {
					status = http.StatusForbidden
				}
				afficherCloture(c, d, walletID, status, err.Error())
				return
			}
			repondreErreurCloture(c, err)
			return
		}
		// Post/Redirect/Get
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/wallets/%d/closure", walletID))
	})

	prive.POST("/wallets/:id/settlements/:settlementId/repay", func(c *gin.Context) {
		walletID, ok := lireID(c, "id")
		if !ok {
			return
		}
		settlementID, ok := lireID(c, "settlementId")
		if !ok {
			return
		}
		compte := CompteCourant(c)

		if err := d.Closure.Rembourser(c.Request.Context(), walletID, settlementID, compte.ID); err != nil {
			repondreErreurCloture(c, err)
			return
		}
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/wallets/%d/closure", walletID))
	})
}

// afficherCloture charge les soldes et rend la page.
func afficherCloture(c *gin.Context, d Deps, walletID int64, status int, erreur string) {
	compte := CompteCourant(c)

	// TODO 1.1 : demander au service les données de la page, pour ce
	// porte-monnaie et ce compte. En cas d'erreur, la traduire avec
	// repondreErreurCloture et s'arrêter là.

	// TODO 1.2 : rendre le template "closure.html" avec le Renderer.
	//
	// Il attend trois clés : "Title", "Page" (les données du service) et
	// "Erreur" (le message à afficher, ou une chaîne vide).
	//
	// Regarde afficherDepenses dans handlers_expenses.go : c'est le même
	// motif, tu peux t'en inspirer.

	_ = compte
	_ = status
	_ = erreur
}

// estRefusCloture : refus « normaux », affichés dans la page.
func estRefusCloture(err error) bool {
	return errors.Is(err, service.ErrClotureReservee) ||
		errors.Is(err, service.ErrDejaClos)
}

// repondreErreurCloture traduit une erreur métier en code HTTP.
//
// Gin ne fait pas cette traduction : une erreur remontée du service est une
// valeur Go ordinaire, pas une exception que le framework intercepterait.
// C'est à nous de décider quel statut elle mérite.
func repondreErreurCloture(c *gin.Context, err error) {
	// TODO 1.3 : associer chaque erreur du paquet service au bon code HTTP.
	//
	// Les erreurs possibles sont déclarées en haut de service/closure.go et
	// service/expense.go, avec un nom qui dit ce qu'elles signifient.
	//
	// Trois familles, trois statuts :
	//   - ce qui n'existe pas       → 404 Not Found
	//   - ce qui est interdit       → 403 Forbidden
	//   - ce qui entre en conflit
	//     avec l'état actuel        → 409 Conflict
	//
	// Indice : errors.Is(err, service.ErrXxx) compare une erreur à une autre,
	// même enveloppée en remontant les couches. Un switch sans condition
	// permet d'enchaîner les cas lisiblement.
	//
	// Le cas par défaut, déjà écrit ci-dessous, couvre les erreurs techniques
	// imprévues : on les journalise sans en montrer le détail au visiteur.

	_ = c.Error(err)
	c.String(http.StatusInternalServerError, "une erreur est survenue, réessayez")
}
