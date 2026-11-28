# Avant le TP Gin — à faire chez vous

Atelier Infrastructures & Frameworks · Groupe Gin · séance du 7 décembre 2026

Comptez **vingt minutes**. Tout se fait sur votre machine : aucun compte à
créer, aucune connexion au réseau de l'université, rien à rendre.

Le TP commence par du code, pas par une installation. Si vous arrivez avec
cette fiche terminée, vous codez dès la première minute.

---

## 1 · Installer trois outils

**Go, version 1.27 ou plus récente** — <https://go.dev/dl>  
Prenez le paquet officiel de votre système : il se charge lui-même d'ajouter
`go` au `PATH`. Fermez et rouvrez votre terminal après l'installation.

**Docker Desktop** — <https://docs.docker.com/get-started>  
La base de données du TP tourne dedans. Docker devra être **lancé** pendant
la séance : sur macOS, la baleine doit être visible dans la barre de menus.

**Un éditeur** — Visual Studio Code, IntelliJ IDEA ou celui que vous avez
l'habitude d'utiliser.

Vérifiez les trois dans un terminal :

```
go version              → go version go1.27… ou plus
docker compose version  → Docker Compose version v2…
git --version           → git version 2…
```

Si une commande n'est pas reconnue, c'est qu'elle n'est pas dans le `PATH` :
réinstallez l'outil, puis rouvrez le terminal.

---

## 2 · Préparer l'éditeur

Go impose son formatage : un fichier mal indenté est considéré comme une
erreur, pas comme une préférence. Configurez le formatage automatique, vous
n'y penserez plus.

**Visual Studio Code** — installez l'extension **Go** officielle (éditeur
Google, identifiant `golang.go`). À la première ouverture d'un fichier `.go`,
elle propose d'installer ses outils : acceptez. Le formatage est actif par
défaut.

**IntelliJ IDEA** — installez le plugin **Go** (Settings → Plugins), puis
Settings → Tools → Actions on Save → cochez *Reformat code*. Attention, ce
plugin demande la version **Ultimate**, gratuite avec votre adresse
universitaire sur <https://www.jetbrains.com/community/education>.

**Un autre éditeur** — vérifiez simplement que `gofmt` est lancé à la
sauvegarde, ou lancez `make fmt` de temps en temps.

---

## 3 · Récupérer le projet

Vous n'avez rien à rendre, donc le plus simple est de cloner directement :

```
git clone https://github.com/miage-if-gin/bons-comptes-tp.git
cd bons-comptes-tp
go mod download
```

Si vous voulez **garder votre travail** après la séance, faites plutôt un
fork : bouton *Fork* en haut à droite de
<https://github.com/miage-if-gin/bons-comptes-tp>, puis clonez votre copie à
la place de la nôtre. Vous pourrez y pousser vos commits.

Ouvrez ensuite le dossier `bons-comptes-tp` dans votre éditeur.

---

## 4 · Créer le fichier d'environnement

L'application ne contient aucun réglage en dur : adresse de la base, port
d'écoute et clé de session viennent tous d'un fichier `.env`. Il n'est pas
dans le dépôt, puisqu'il contiendrait des secrets en production. Vous le
créez à partir du modèle fourni :

```
cp .env.example .env            (macOS, Linux, Git Bash)
Copy-Item .env.example .env     (PowerShell)
```

**Il n'y a rien à compléter.** Les valeurs par défaut pointent vers la base
locale que Docker va lancer.

---

## 5 · Lancer l'application

Docker Desktop doit tourner. Puis, à la racine du projet :

```
make up      démarre la base, attendez l'affichage « healthy »
make seed    remplit la base : 40 comptes, 25 porte-monnaies, 5 000 dépenses
make run     démarre le serveur sur http://localhost:8080
```

**Sous Windows**, `make` n'existe pas. Le plus simple est de travailler dans
**WSL 2**, où tout ce qui précède fonctionne tel quel. Sinon, remplacez les
commandes par : `docker compose up -d`, `go run ./cmd/seed`,
`go run ./cmd/server` — et pensez à définir vous-même les variables du
`.env`, puisque c'est le `Makefile` qui s'en charge d'habitude.

---

## 6 · Vérifier que tout est en place

Deux vérifications, l'une dans le navigateur, l'autre dans le terminal.

**Dans le navigateur.** Il vous faut l'adresse d'un compte. Ouvrez un second
terminal et lancez :

```
docker exec -it bons-comptes-tp-db psql -U bonscomptes -d bonscomptes -c "SELECT email FROM accounts LIMIT 3;"
```

Allez sur <http://localhost:8080>, connectez-vous avec l'une de ces adresses
et le mot de passe `demo1234`. Ouvrez **Mes porte-monnaies**, choisissez-en
un, puis **Soldes et clôture**.

> La page doit être **blanche**. C'est le point de départ du TP : le code qui
> la remplit est précisément ce que vous allez écrire.

**Dans le terminal.** Arrêtez le serveur avec `Ctrl+C`, puis :

```
make point
```

Cinq lignes doivent s'afficher, **toutes rouges**. C'est le tableau de bord
que vous utiliserez pendant les deux heures.

---

## Si quelque chose ne marche pas

Ne restez pas bloqué et ne réinstallez pas tout : **notez le message d'erreur
exact** et envoyez-le-nous, ou montrez-le-nous en début de séance. Dix
minutes sont prévues pour ça.

Deux cas fréquents, que vous pouvez régler seul :

**Le port 8080 est déjà pris** (`bind: address already in use`) — un autre
serveur tourne. Ouvrez le fichier `.env` et remplacez `PORT=8080` par
`PORT=8090`, puis relancez `make run`. L'application sera alors sur
<http://localhost:8090>.

**Docker refuse de démarrer** (`Cannot connect to the Docker daemon`) —
Docker Desktop n'est pas lancé. Ouvrez-le, attendez qu'il soit prêt, relancez
`make up`.

---

Thomas Devinck · Yuliia Los · Sofia Demchuk · Baptiste Delisle
