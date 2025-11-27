-Slypot13
-gabpr278

# Projet Power'4 Web 🔴🟡

Bonjour ! Voici notre projet de **Puissance 4** développé en Go pour le module de développement web.

[cite_start]Le but était de créer un jeu fonctionnel jouable à deux sur le même ordinateur, avec une interface web[cite: 4, 5].

## 👥 L'équipe (Binôme)
* **[TON PRÉNOM]**
* **[PRÉNOM DE TON BINÔME]**

## 🚀 Comment lancer le projet

Pour tester notre jeu, suivez ces étapes simple :

1.  Clonez le dépôt :
    ```bash
    git clone [LIEN_DU_REPO_GITHUB]
    ```
2.  Ouvrez le dossier dans votre terminal.
3.  Lancez le serveur :
    ```bash
    go run main.go
    ```
4.  Ouvrez votre navigateur et allez à l'adresse :
    `http://localhost:8080` (ou le port indiqué dans le terminal).

## 🎮 Fonctionnalités du jeu

Nous avons respecté les règles classiques du Puissance 4 et les consignes du sujet :

* [cite_start]**Mode multijoueur local** : Deux joueurs s'affrontent tour par tour sur la même machine[cite: 5, 86].
* [cite_start]**Grille** : Une grille de 7 colonnes et 6 lignes[cite: 13].
* [cite_start]**Victoire** : Le jeu détecte si un joueur aligne 4 jetons (ligne, colonne ou diagonale)[cite: 18].
* [cite_start]**Égalité** : Le jeu détecte si la grille est pleine sans vainqueur[cite: 19].
* [cite_start]**Historique** : Une page Scoreboard permet de voir les résultats des parties précédentes[cite: 55].

## 📂 Structure et Pages

[cite_start]Le site est organisé avec les routes demandées[cite: 88]:

* [cite_start]`/` : **Accueil** - Présentation du projet et règles du jeu[cite: 28].
* [cite_start]`/game/init` : **Initialisation** - Choix des pseudos et des couleurs des jetons[cite: 34].
* `/game/play` : **Jeu** - La grille où on joue. [cite_start]On a utilisé des formulaires pour choisir les colonnes[cite: 42, 45].
* [cite_start]`/game/end` : **Fin** - Affiche le gagnant et un bouton pour rejouer[cite: 49].
* [cite_start]`/game/scoreboard` : **Scores** - Historique des parties jouées (sauvegardé en mémoire)[cite: 55].

## 🛠️ Technologies utilisées

[cite_start]Comme demandé dans les contraintes techniques, nous n'avons pas utilisé de framework externe[cite: 73, 75].

* [cite_start]**Langage** : Golang (Go) pour le serveur[cite: 77].
* [cite_start]**Frontend** : HTML, CSS et Templates GOHTML[cite: 78].
* [cite_start]**JavaScript** : Uniquement pour quelques petites animations légères[cite: 74].

---
*Projet réalisé pour Ynov Campus Aix - 2025*