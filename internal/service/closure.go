package service

import (
	"context"
	"errors"

	"github.com/miage-if-gin/bons-comptes-tp/internal/domain"
	"github.com/miage-if-gin/bons-comptes-tp/internal/repository"
)

var (
	ErrClotureReservee       = errors.New("seul le créateur du porte-monnaie peut le clôturer")
	ErrDejaClos              = errors.New("ce porte-monnaie est déjà clos")
	ErrPasEncoreClos         = errors.New("le porte-monnaie doit être clos avant de rembourser")
	ErrVirementIntrouvable   = errors.New("virement introuvable")
	ErrRemboursementInterdit = errors.New("seul le compte qui doit l'argent peut effectuer ce remboursement")
	ErrDejaRembourse         = errors.New("ce virement a déjà été remboursé")
)

type Solde struct {
	AccountID int64
	Nom       string
	PayeCents int64
	PartCents int64
}

func (s Solde) SoldeCents() int64 { return s.PayeCents - s.PartCents }

func (s Solde) Paye() string { return FormatMontant(s.PayeCents) }
func (s Solde) Part() string { return FormatMontant(s.PartCents) }
func (s Solde) Net() string  { return FormatMontant(s.SoldeCents()) }

func (s Solde) Montant() string {
	if s.SoldeCents() < 0 {
		return FormatMontant(-s.SoldeCents())
	}
	return FormatMontant(s.SoldeCents())
}

type Virement struct {
	De           int64
	DeNom        string
	Vers         int64
	VersNom      string
	MontantCents int64
}

func (v Virement) Montant() string { return FormatMontant(v.MontantCents) }

// CalculerSoldes établit, pour chaque membre, ce qu'il a payé, ce qu'il
// devait payer, et donc ce qu'il doit ou ce qu'on lui doit.
//
// La règle : le total des dépenses attribuées se partage à parts égales
// entre tous les membres acceptés. Quelqu'un qui a payé plus que sa part
// est créancier, quelqu'un qui a moins payé est débiteur.
//
// Deux cas à ne pas manquer :
//   - la division ne tombe pas toujours juste (100 € entre 3) ;
//   - un payeur peut ne plus être membre accepté : il n'a pas de part à
//     payer, mais on lui doit quand même ce qu'il a avancé.
func CalculerSoldes(membres []domain.Account, payes []domain.AccountPaid) []Solde {
	// TODO 3 : construire la liste des soldes.
	//
	// Le test internal/service/closure_test.go est la spécification. Lis-le
	// avant de coder : il dit ce qui est attendu dans chaque cas, y compris
	// les cas limites.
	//
	// Invariant à respecter : la somme de tous les SoldeCents() doit valoir
	// exactement zéro. Pas « à peu près » — exactement.
	//
	// Note : Go refuse de compiler avec un import inutilisé. Les paquets dont
	// tu auras besoin (sort, par exemple) sont donc à ajouter toi-même —
	// IntelliJ et VS Code le font automatiquement à la sauvegarde.

	return nil
}

// MinimiserVirements transforme une liste de soldes en une liste de
// virements qui remet tout le monde à zéro, en limitant leur nombre.
//
// Chercher le minimum absolu de virements est un problème NP-difficile : il
// faudrait tester tous les sous-groupes dont les soldes s'annulent. On se
// contente d'une heuristique, qui garantit au plus n − 1 virements pour
// n comptes — largement suffisant en pratique.
//
// Trois invariants, vérifiés par le test :
//   - tout montant est strictement positif ;
//   - personne ne se vire de l'argent à soi-même ;
//   - après tous les virements, chaque solde vaut exactement zéro.
func MinimiserVirements(soldes []Solde) []Virement {
	// TODO 4 : produire les virements.
	//
	// Piste : séparer les débiteurs (solde négatif) des créanciers (solde
	// positif), puis apparier les uns aux autres jusqu'à épuisement.
	//
	// Le test internal/service/closure_test.go est la spécification.
	// TestMinimiserVirements_Aleatoire vérifie les invariants sur 500 cas
	// tirés au hasard : s'il passe, l'algorithme est correct.

	return nil
}

// Service : vérifie les droits, lit la base, appelle les calculs purs.

// LigneVirement est un virement enregistré  tel qu'affiché après clôture.
type LigneVirement struct {
	domain.SettlementLine
	PeutRembourser bool
}

func (l LigneVirement) Montant() string { return FormatMontant(l.AmountCents) }
func (l LigneVirement) Rembourse() bool { return l.SettledAt != nil }

// PageCloture rassemble tout ce qu'affiche la page de clôture.
type PageCloture struct {
	Wallet            domain.Wallet
	Soldes            []Solde
	Apercu            []Virement      // porte-monnaie ouvert : ce que donnerait la clôture
	Virements         []LigneVirement // porte-monnaie clos : virements enregistrés
	TotalCents        int64
	NonAttribuesCents int64
	EstProprietaire   bool
}

func (p PageCloture) Ouvert() bool  { return p.Wallet.Status == "open" }
func (p PageCloture) Total() string { return FormatMontant(p.TotalCents) }

// NonAttribues : total des dépenses sans payeur, affiché à part.
func (p PageCloture) NonAttribues() string { return FormatMontant(p.NonAttribuesCents) }

// Debiteurs : la liste « qui doit donner de l'argent, et combien ».
func (p PageCloture) Debiteurs() []Solde {
	var res []Solde
	for _, s := range p.Soldes {
		if s.SoldeCents() < 0 {
			res = append(res, s)
		}
	}
	return res
}

// Crediteurs : la liste « qui doit recevoir de l'argent, et combien ».
func (p PageCloture) Crediteurs() []Solde {
	var res []Solde
	for _, s := range p.Soldes {
		if s.SoldeCents() > 0 {
			res = append(res, s)
		}
	}
	return res
}

type ClosureService struct {
	wallets     *repository.WalletRepository
	settlements *repository.SettlementRepository
	expenses    *repository.ExpenseRepository
}

func NewClosureService(w *repository.WalletRepository, s *repository.SettlementRepository, e *repository.ExpenseRepository) *ClosureService {
	return &ClosureService{wallets: w, settlements: s, expenses: e}
}

// Consulter renvoie les soldes du porte-monnaie, avec l'aperçu des virements
func (s *ClosureService) Consulter(ctx context.Context, walletID, accountID int64) (*PageCloture, error) {
	wallet, err := s.verifierAcces(ctx, walletID, accountID)
	if err != nil {
		return nil, err
	}

	soldes, err := s.soldes(ctx, walletID)
	if err != nil {
		return nil, err
	}

	page := &PageCloture{
		Wallet:          *wallet,
		Soldes:          soldes,
		EstProprietaire: wallet.OwnerID == accountID,
	}
	// Le total vient de la table des dépenses, comme sur la page des
	// dépenses : il inclut les dépenses sans payeur.
	_, total, err := s.expenses.Totals(ctx, walletID)
	if err != nil {
		return nil, err
	}
	page.TotalCents = total

	// Ce qui n'a été avancé par personne n'entre pas dans les soldes.
	page.NonAttribuesCents = total
	for _, so := range soldes {
		page.NonAttribuesCents -= so.PayeCents
	}
	if page.Ouvert() {
		page.Apercu = MinimiserVirements(soldes)
		return page, nil
	}

	lignes, err := s.settlements.ListByWallet(ctx, walletID)
	if err != nil {
		return nil, err
	}
	for _, l := range lignes {
		page.Virements = append(page.Virements, LigneVirement{
			SettlementLine: l,
			PeutRembourser: l.SettledAt == nil && l.FromAccountID == accountID,
		})
	}
	return page, nil
}

// Cloturer fige le porte-monnaie et enregistre les virements à effectuer
func (s *ClosureService) Cloturer(ctx context.Context, walletID, accountID int64) error {
	wallet, err := s.verifierAcces(ctx, walletID, accountID)
	if err != nil {
		return err
	}
	if wallet.OwnerID != accountID {
		return ErrClotureReservee
	}
	if wallet.Status != "open" {
		return ErrDejaClos
	}

	soldes, err := s.soldes(ctx, walletID)
	if err != nil {
		return err
	}

	var aEnregistrer []domain.Settlement
	for _, v := range MinimiserVirements(soldes) {
		aEnregistrer = append(aEnregistrer, domain.Settlement{
			WalletID:      walletID,
			FromAccountID: v.De,
			ToAccountID:   v.Vers,
			AmountCents:   v.MontantCents,
		})
	}

	err = s.settlements.Close(ctx, walletID, aEnregistrer)
	if errors.Is(err, repository.ErrWalletAlreadyClosed) {
		return ErrDejaClos
	}
	return err
}

func (s *ClosureService) Rembourser(ctx context.Context, walletID, settlementID, accountID int64) error {
	wallet, err := s.verifierAcces(ctx, walletID, accountID)
	if err != nil {
		return err
	}
	if wallet.Status != "closed" {
		return ErrPasEncoreClos
	}

	virement, err := s.settlements.FindInWallet(ctx, walletID, settlementID)
	if err != nil {
		return err
	}
	if virement == nil {
		return ErrVirementIntrouvable
	}
	if virement.FromAccountID != accountID {
		return ErrRemboursementInterdit
	}
	if virement.SettledAt != nil {
		return ErrDejaRembourse
	}
	return s.settlements.MarkSettled(ctx, walletID, settlementID)
}

func (s *ClosureService) soldes(ctx context.Context, walletID int64) ([]Solde, error) {
	membres, err := s.wallets.AcceptedMembers(ctx, walletID)
	if err != nil {
		return nil, err
	}
	payes, err := s.settlements.PaidByAccount(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return CalculerSoldes(membres, payes), nil
}

// verifierAcces renvoie le porte-monnaie si le compte en est membre accepté.
// Mêmes erreurs que pour les dépenses : le handler les traduit pareil.
func (s *ClosureService) verifierAcces(ctx context.Context, walletID, accountID int64) (*domain.Wallet, error) {
	wallet, err := s.wallets.FindByID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrPorteMonnaieIntrouvable
	}

	membre, err := s.wallets.IsAcceptedMember(ctx, walletID, accountID)
	if err != nil {
		return nil, err
	}
	if !membre {
		return nil, ErrAccesRefuse
	}
	return wallet, nil
}
