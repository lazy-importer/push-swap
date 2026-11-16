Projet Go pour le sujet **Push-Swap** : deux exécutables, `push-swap` et
`checker`, avec les onze instructions du sujet.

---

## 1. Compilation et lancement

```console
$ go build -o push-swap ./cmd/push-swap
$ go build -o checker ./cmd/checker
```

```console
$ ./push-swap "4 67 3 87 23"
$ ./push-swap "4 67 3 87 23" | ./checker "4 67 3 87 23"
OK
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
├── go.mod
├── cmd/
│   ├── push-swap/main.go     # point d'entrée du générateur d'instructions
│   └── checker/main.go       # point d'entrée du vérificateur
└── internal/pushswap/
    ├── stack.go              # piles, 11 instructions, rotation la plus courte
    ├── parse.go              # arguments et doublons
    ├── sort.go               # tri : trois éléments + insertion au moindre coût
    ├── checker.go            # lecture et exécution des instructions
    └── *_test.go             # tests unitaires
```

La séparation suit la grille : `parse.go` (lecture et validation),
`stack.go` (opérations centralisées), `sort.go` (algorithme),
`checker.go` (vérification), `cmd/` (interface des deux binaires).

## 4. Algorithme

Une seule stratégie, pour toutes les tailles :

1. On pousse les éléments de `a` sur `b` jusqu'à ne garder que trois éléments.
2. On trie ces trois éléments (2 instructions au plus).
3. Pour chaque élément de `b`, on calcule le coût de son insertion dans `a` :
   rotation de `b` + rotation de `a` + `pa`. L'élément le moins coûteux est
   inséré. `a` reste triée circulairement ; la cible est trouvée par recherche
   dichotomique.
4. Quand les deux rotations vont dans le même sens, `rr` ou `rrr` les exécute
   ensemble, ce qui économise des instructions.
5. Un dernier quart de tour amène le minimum en haut de `a`.

C'est une heuristique : elle ne garantit pas le minimum absolu d'instructions,
mais reste raisonnable. Mesures (piles aléatoires reproductibles) :

| n | Instructions (moyenne) | Vérification |
|---|---|---|
| 100 | ~590 | OK |
| 500 | ~5 200 | OK |
| 1 000 | ~14 200 | OK |

## 5. Tests

```console
$ go test ./...          # suite complète
$ go test -short ./...   # sans les tests exhaustifs
```

Les tests couvrent :

- les onze instructions, y compris `ss`, `rr` et `rrr`, et les opérations sur
  des piles trop petites ;
- la lecture des arguments, les nombres négatifs, les décimales, les
  débordements et les doublons ;
- le `checker` : `OK`, `KO`, instructions inconnues, espaces, lignes vides,
  fins de ligne Windows ;
- le tri de toutes les permutations jusqu'à 3 éléments, et des piles
  aléatoires de 8 à 250 éléments, rejouées par le `checker` ;
- les exemples du sujet, une pile déjà triée (0 instruction).

## 6. Correspondance avec la grille d'évaluation

| Critère | Où / comment |
|---|---|
| **C1** Arguments et initialisation | `parse.go`, `cmd/push-swap/main.go` ; premier argument = sommet ; `b` vide ; aucun argument = aucune sortie |
| **C2** Instructions Push-Swap | `stack.go` : `pa`, `pb`, `sa`, `sb`, `ss`, `ra`, `rb`, `rr`, `rra`, `rrb`, `rrr` |
| **C3** `push-swap` fonctionnel | `Sort()` vérifie `Done()` ; seules les instructions sont écrites |
| **C4** Optimisation | insertion au moindre coût avec `rr`/`rrr` combinés, rotation la plus courte ; mesures ci-dessus |
| **C5** `checker` | `checker.go` + `cmd/checker/main.go` ; lecture stdin, `OK`/`KO` |
| **C6** Gestion des erreurs | `ParseArgs`, `CheckDuplicates`, `Apply` ; `Error` + code 1 |
| **C7** Contraintes | Go, bibliothèque standard uniquement, binaires `push-swap` et `checker` |
| **C8** Algorithmique et piles | `Stack` centralise les opérations ; `sort.go` documenté |
| **C9** Architecture | `cmd/` + `internal/pushswap/`, responsabilités séparées |
| **C10** Qualité du code | fonctions courtes, code commenté, `gofmt`, `go vet`, tests |
| **C11** Robustesse | cas limites testés (vide, un élément, déjà trié, doublons, entrée invalide) |

## 7. Limites connues

- L'insertion est une heuristique : elle ne garantit pas le minimum absolu
  d'instructions (par exemple 7 instructions pour `4 67 3 87 23`, là où une
  recherche exhaustive trouverait 6).
- `ss`, `sb`, `rrb` et `rrr` ne sont pas produits par le tri, mais ils sont
  implémentés et testés : le `checker` doit accepter tout le langage du sujet.
