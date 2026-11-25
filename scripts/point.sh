#!/usr/bin/env bash
# Où j'en suis : une ligne par mission du TP.
# Lancer avec « make point » : c'est le Makefile qui charge le .env.
set -u

vert=$'\033[32m'; rouge=$'\033[31m'; gris=$'\033[90m'; zero=$'\033[0m'

mission() {
  local nom="$1" paquet="$2" motif="$3"
  if go test "$paquet" -count=1 -run "$motif" >/dev/null 2>&1; then
    printf '   %s✓%s  %s\n' "$vert" "$zero" "$nom"
  else
    printf '   %s✗%s  %s\n' "$rouge" "$zero" "$nom"
  fi
}

echo
echo "   La clôture des bons comptes"
echo
mission "Mission 1 · le handler et les codes HTTP" ./internal/http/ 'Cloture'
mission "Mission 2 · qui a payé combien"           ./internal/repository/ 'PaidByAccount'
mission "Mission 3 · les soldes"                   ./internal/service/ 'CalculerSoldes'
mission "Mission 4 · les virements"                ./internal/service/ 'MinimiserVirements_(CasSimple|Aleatoire)'
echo
mission "Bonus A · les paires exactes"             ./internal/service/ 'MinimiserVirements_PairesExactes'
echo
printf '   %sLe détail d une mission : make m1 · m2 · m3 · m4%s\n' "$gris" "$zero"
printf '   %sLes missions 1 et 2 ont besoin de la base : make up%s\n' "$gris" "$zero"
echo
