# Push-Swap — livrable d'évaluation

Projet Go complet pour le sujet **Push-Swap** (`Push Swap - Sujet Ytrack (1).md`)
et la grille `Push swap - Grille d'évaluation groupe (1).xlsx`.

Le code des séances a servi de point de départ : structure `Piles`/`Stack`, les
huit opérations de base, le simulateur, la lecture des arguments, la validation
des doublons, le calcul des rangs et le tri à trois éléments. Le tout a été
complété, corrigé et étendu pour produire les deux exécutables demandés
(`push-swap` et `checker`) avec les onze instructions du sujet.

---

## 1. Compilation et lancement

```console
$ make build          # produit ./push-swap et ./checker
$ make test           # go test ./...
$ make bench          # benchmarks des stratégies
```

Sans `make` :

```console
$ go build -o push-swap ./cmd/push-swap
$ go build -o checker ./cmd/checker
```

Les deux exécutables attendus sont à la racine de `evals/` :

```console
$ ./push-swap "4 67 3 87 23"
$ ./push-swap "4 67 3 87 23" | ./checker "4 67 3 87 23"
OK
$ ARG="4 67 3 87 23"; ./push-swap "$ARG" | wc -l
6
```

`push-swap` accepte aussi un argument par nombre :

```console
$ ./push-swap 2 1 3 6 5 8
pb
pb
ra
sa
rrr
pa
pa
```

## 2. Comportement attendu

| Situation | Résultat |
|---|---|
| Aucun argument | aucune sortie, code 0 |
| Arguments valides | opérations seules, une par ligne, code 0 |
| Pile déjà triée | aucune sortie, code 0 |
| Argument non entier | `Error` sur la sortie d'erreur, code 1 |
| Doublon | `Error` sur la sortie d'erreur, code 1 |
| `checker` : `a` triée et `b` vide | `OK` |
| `checker` : sinon | `KO` |
| `checker` : instruction inconnue ou mal formée | `Error` sur la sortie d'erreur, code 1 |
| `checker` sans argument | aucune sortie, code 0 |

`push-swap` n'écrit **jamais** autre chose que les instructions sur la sortie
standard : pas de titre, pas de compteur, pas d'état des piles.

## 3. Arborescence

```text
evals/
├── Makefile
├── go.mod
├── cmd/
│   ├── push-swap/main.go     # point d'entrée du générateur d'instructions
│   └── checker/main.go       # point d'entrée du vérificateur
└── internal/pushswap/
    ├── stack.go              # piles, 11 instructions, rotation la plus courte
    ├── parse.go              # arguments, doublons, rangs
    ├── sort.go              # stratégies de tri
    ├── optimal.go            # recherche exacte pour les petites piles
    ├── checker.go            # lecture et exécution des instructions
    └── *_test.go             # tests unitaires et exhaustifs
```

La séparation suit la grille : `parse.go` (lecture et validation),
`stack.go` (opérations centralisées), `sort.go`/`optimal.go` (algorithmes),
`checker.go` (vérification), `cmd/` (interface des deux binaires).

## 4. Algorithmes

### 4.1 Rangs

Les nombres reçus sont remplacés par leurs **rangs** : `0` pour le plus petit,
`n-1` pour le plus grand. Les rangs conservent l'ordre, donc la suite
d'instructions est identique, mais toutes les valeurs tiennent entre `0` et
`n-1`. Exemple : `[40, 10, 50, 20, 30]` devient `[3, 0, 4, 1, 2]`.

### 4.2 Petites piles : recherche exacte (n ≤ 7)

Pour 2 à 7 éléments, le programme fait un **parcours en largeur** sur l'espace
complet des états `(a, b)` avec les onze instructions. Le premier état trié
atteint est donc à distance minimale : la solution est **prouvablement la plus
courte**.

L'espace reste petit : 720 états pour 5 éléments, 5 040 pour 6, 40 320 pour 7.
Coût observé : environ 0,5 ms pour 5 éléments, 90 ms au pire pour 7.

Totaux exacts vérifiés par les tests sur **toutes** les permutations :

| n | Permutations | Total optimal d'instructions | Maximum par pile |
|---|---|---|---|
| 2 | 2 | 1 | 1 |
| 3 | 6 | 7 | 2 |
| 4 | 24 | 69 | 5 |
| 5 | 120 | 583 | 8 |
| 6 | 720 | 5 083 | 10 |
| 7 | 5 040 | 46 903 | 13 |

### 4.3 Piles moyennes : insertion au meilleur coût (8 ≤ n ≤ 1 000)

1. On pousse les éléments de `a` sur `b` jusqu'à ne garder que trois éléments.
2. On trie ces trois éléments (2 instructions au plus).
3. Pour chaque élément de `b`, on calcule le coût de son insertion dans `a` :
   rotation de `b` + rotation de `a` + `pa`. L'élément le moins coûteux est
   inséré. `a` reste triée circulairement ; la cible est trouvée par recherche
   dichotomique.
4. Quand les deux rotations vont dans le même sens, `rr` ou `rrr` les exécute
   ensemble, ce qui économise des instructions.
5. Un dernier quart de tour amène le minimum en haut de `a`.

### 4.4 Très grandes piles : stratégie par blocs (n > 1 000)

Les rangs sont découpés en environ `√n` blocs, poussés sur `b` du plus petit au
plus grand, avec `rb` pour ranger la moitié basse de chaque bloc au fond. Les
plus grandes valeurs se retrouvent près du sommet de `b` ; il suffit alors de
sortir le maximum par le chemin le plus court avec `pa`. Cette stratégie a une
meilleure croissance asymptotique que l'insertion, au prix de plus de
mouvements sur les tailles usuelles.

### 4.5 Mesures (20 piles aléatoires par taille, mêmes entrées pour `checker`)

| n | Moyenne d'instructions | Minimum | Maximum | Vérification |
|---|---|---|---|---|
| 5 | 4,8 | 0 | 8 | 20/20 OK |
| 6 | 7,6 | 6 | 10 | 20/20 OK |
| 7 | 9,1 | 6 | 11 | 20/20 OK |
| 10 | 26,1 | 22 | 31 | 20/20 OK |
| 25 | 93,0 | 80 | 105 | 20/20 OK |
| 50 | 231,6 | 211 | 254 | 20/20 OK |
| 100 | 582,5 | 538 | 622 | 20/20 OK |
| 500 | 5 216,9 | 5 027 | 5 426 | 20/20 OK |
| 1 000 | 13 715,0 | 13 347 | 14 482 | 20/20 OK |

## 5. Tests

```console
$ go test ./...          # suite complète (environ 8 s)
$ go test -short ./...   # sans les recherches exhaustives
```

Les tests couvrent :

- les onze instructions, y compris `ss`, `rr` et `rrr`, et les opérations sur
  des piles trop petites ;
- la lecture des arguments, les nombres négatifs, les décimales, les
  débordements, les doublons et les rangs ;
- le `checker` : `OK`, `KO`, instructions inconnues, espaces, lignes vides,
  fins de ligne Windows ;
- le tri de **toutes** les permutations de 0 à 6 éléments, avec vérification du
  nombre total optimal d'instructions ;
- un échantillon de piles de 7 éléments (borne de 13 instructions) ;
- des piles aléatoires de 8 à 500 éléments, rejouées par le `checker` ;
- l'égalité stricte des instructions obtenues avec les nombres d'origine et
  avec leurs rangs ;
- des seuils de qualité sur les tailles 100 et 500.

## 6. Correspondance avec la grille d'évaluation

| Critère | Où / comment |
|---|---|
| **C1** Arguments et initialisation | `parse.go`, `cmd/push-swap/main.go` ; premier argument = sommet ; `b` vide ; aucun argument = aucune sortie |
| **C2** Instructions Push-Swap | `stack.go` : `pa`, `pb`, `sa`, `sb`, `ss`, `ra`, `rb`, `rr`, `rra`, `rrb`, `rrr` |
| **C3** `push-swap` fonctionnel | `Sort()` vérifie `Done()` ; seules les instructions sont écrites |
| **C4** Optimisation | recherche exacte ≤ 7, insertion au meilleur coût avec `rr`/`rrr`, blocs au-delà ; mesures ci-dessus |
| **C5** `checker` | `checker.go` + `cmd/checker/main.go` ; lecture stdin, `OK`/`KO` |
| **C6** Gestion des erreurs | `ParseArgs`, `CheckDuplicates`, `Apply` ; `Error` + code 1 |
| **C7** Contraintes | Go, bibliothèque standard uniquement, binaires `checker` et `push-swap` via `Makefile` |
| **C8** Algorithmique et piles | `Stack` centralise les opérations ; `sort.go` et `optimal.go` documentés |
| **C9** Architecture | `cmd/` + `internal/pushswap/`, responsabilités séparées |
| **C10** Qualité du code | fonctions courtes, code commenté, `gofmt`, `go vet`, tests |
| **C11** Robustesse | cas limites testés (vide, un élément, déjà trié, doublons, entrée invalide, combinaisons) |

## 7. Limites connues

- La recherche exacte est limitée à 7 éléments : au-delà, l'espace d'états
  croît en `(n+1)!` et devient trop grand.
- L'insertion est une heuristique : elle ne garantit pas l'optimum global,
  mais reste nettement sous les seuils habituels (100 et 500 éléments).
- `ss`, `rr` et `rrr` ne sont pas produits par les stratégies générales, sauf
  `rr`/`rrr` dans l'insertion ; ils sont néanmoins implémentés et testés, car
  le `checker` doit accepter tout le langage du sujet.
