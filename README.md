# Les bons comptes — TP Gin

Squelette du TP de l'atelier **Gin (Go)**, module Infrastructures & Frameworks,
M2 MIAGE IFI — séance du 7 décembre 2026.

L'application gère des dépenses partagées : on crée un porte-monnaie, on y
invite des comptes, chacun déclare ses dépenses. **La clôture — calculer qui
doit combien à qui — est à écrire.** C'est l'objet du TP.

---

## Avant la séance

| Prérequis | Version | Vérifier |
|---|---|---|
| [Go](https://go.dev/dl) | ≥ 1.27 | `go version` |
| [Docker + Compose](https://docs.docker.com/get-started) | — | `docker compose version` |

```bash
git clone https://github.com/miage-if-gin/bons-comptes-tp.git
cd bons-comptes-tp
go mod download

cp .env.example .env
make up                       # PostgreSQL, attendre « healthy »
make seed                     # 40 comptes, 25 porte-monnaies, 5 000 dépenses
make run                      # http://localhost:8080
```

Connecte-toi avec n'importe quel compte généré — mot de passe `demo1234`.
Pour en trouver un :

```bash
docker exec -it bons-comptes-tp-db psql -U bonscomptes -d bonscomptes \
  -c "SELECT email FROM accounts LIMIT 3;"
```

Ouvre un porte-monnaie, puis **Soldes et clôture** : la page est vide. À toi
de la remplir.

### Si quelque chose ne démarre pas

**Le port 5432 ou 8080 est déjà pris** — un autre PostgreSQL ou un autre
serveur tourne. Arrête-le, ou change le port dans `docker-compose.yml` et
dans le `.env`.

**Sur Windows**, `make` n'existe pas. Travaille dans WSL 2, ou remplace les
cibles par leurs commandes :

| Cible | Commande |
|---|---|
| `make up` | `docker compose up -d` |
| `make run` | `go run ./cmd/server` |
| `make seed` | `go run ./cmd/seed` |
| `make test` | `go test ./... -race -cover` |

Sans `make`, il faut exporter les variables du `.env` à la main : c'est le
`Makefile` qui les charge.

---

## Le travail

Quatre fonctions à écrire, dans cet ordre. Chacune rend la page de clôture
un peu plus complète.

```bash
grep -rn "TODO" --include="*.go" .
```

| | Fichier | Quoi |
|---|---|---|
| 1 | `internal/http/handlers_closure.go` | le handler et la traduction des erreurs |
| 2 | `internal/repository/settlement.go` | le SQL qui dit qui a payé combien |
| 3 | `internal/service/closure.go` | le calcul des soldes |
| 4 | `internal/service/closure.go` | la minimisation des virements |

**Les tests sont la spécification.** Ils disent précisément ce qui est
attendu, y compris les cas limites. Lis-les avant de coder.

```bash
make test
```

Le sujet du TP détaille chaque mission.

---

## L'architecture en quatre couches

Gin n'impose aucune structure : celle-ci est un choix, pas une convention du
framework.

```
cmd/server/     point d'entrée : configuration, base, serveur
cmd/seed/       génération du jeu de données
internal/
  config/       configuration lue depuis l'environnement
  db/           connexion PostgreSQL et migrations
  domain/       structures métier, sans dépendance technique
  repository/   accès aux données — le SQL est écrit à la main
  service/      règles métier, testables sans serveur ni base
  http/         routeur, handlers, middlewares
migrations/     schéma SQL versionné, embarqué dans le binaire
web/            templates HTML et fichiers statiques, embarqués aussi
```

Une requête traverse les couches de haut en bas : le routeur appelle un
handler, qui appelle un service, qui appelle un repository. Jamais l'inverse,
et un handler n'écrit jamais de SQL.

| Besoin | Choix | Pourquoi |
|---|---|---|
| Accès données | GORM + SQL brut | aucune installation ; le SQL d'agrégation reste explicite |
| Migrations | golang-migrate | SQL versionné, embarqué via `go:embed` |
| Vues | `html/template` | pas de build frontend, échappement XSS natif |
| CSS | Tailwind (CDN) | aucune installation |
| Sessions | table `sessions` | révocation réelle à la déconnexion |
| Montants | `BIGINT` en centimes | jamais de flottant pour de l'argent |

---

## Commandes

```
make help     liste toutes les commandes
make up       démarre PostgreSQL et Adminer (http://localhost:8081)
make run      lance le serveur
make seed     peuple la base — options via `go run ./cmd/seed -h`
make test     lance les tests
make fmt      formate le code (gofmt est obligatoire en Go)
make psql     shell SQL sur la base
make down     arrête les conteneurs
make reset    détruit la base et ses données
```
