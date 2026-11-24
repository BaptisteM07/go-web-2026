# Lexique

| Terme | Ce que ça veut dire |
| --- | --- |
| HTTP | Le protocole du web : une requête part du client, une réponse revient du serveur. Tout est du texte. |
| Requête / réponse | Ce que le client envoie (méthode, URL, en-têtes) et ce que le serveur renvoie (code, en-têtes, contenu). |
| URL | L'adresse d'une ressource : schéma, hôte, port, chemin, query string. |
| Query string | La partie après le `?` d'une URL : `?genre=course&q=mario`. Lue en Go avec `r.URL.Query().Get("q")`. |
| Méthode | Ce qu'on veut faire : GET lire, POST envoyer, PUT ou PATCH modifier, DELETE supprimer. |
| Status code | Le code de la réponse : 2xx tout va bien, 3xx redirection, 4xx erreur du client, 5xx erreur du serveur. |
| 401 / 403 | 401 : le serveur ne sait pas qui vous êtes (pas de clé, pas connecté). 403 : il sait qui vous êtes, mais vous n'avez pas le droit. |
| OpenAPI / Swagger | OpenAPI : un format standard pour décrire une API (ses endpoints, ses paramètres, ses réponses). Swagger UI : la page web qui affiche cette description et permet d'essayer chaque endpoint. |
| Clé d'API | Un code secret, souvent envoyé dans un en-tête, qui permet à l'API de savoir qui l'appelle. |
| En-tête (header) | Une information sur la requête ou la réponse, par exemple `Content-Type: application/json`. |
| JSON | Un format texte pour échanger des données : objets `{}`, tableaux `[]`, chaînes, nombres, booléens, `null`. |
| API | Un programme qui offre des données ou des actions à d'autres programmes. |
| API REST | Une API web où chaque URL désigne une ressource (`/api/games/3`) et où la méthode dit quoi en faire. |
| Endpoint | Une URL précise d'une API : `GET /api/games` est un endpoint. |
| Backend | La partie qui tourne sur le serveur : ici, le programme Go. |
| Frontend | La partie qui tourne dans le navigateur : HTML, CSS, JavaScript. |
| Handler | En Go, la fonction qui reçoit une requête et écrit la réponse : `func(w http.ResponseWriter, r *http.Request)`. |
| Routeur (ServeMux) | Ce qui choisit le handler en fonction de la méthode et de l'URL. |
| Fichiers statiques | Les fichiers envoyés tels quels : CSS, JavaScript, images. |
| Template | Un fichier HTML avec des trous (`{{.Name}}`) que Go remplit avec des données. |
| Client HTTP | Le programme qui envoie une requête. Notre serveur Go est client de l'API externe. |
| Timeout | La durée maximale qu'on accepte d'attendre une réponse avant d'abandonner. |
| Décoder / encoder | Transformer du JSON en données Go (Decode), ou des données Go en JSON (Encode). |
| Struct tag | L'annotation `json:"name"` qui relie un champ Go à une clé JSON. |
| Cache | Une copie des données gardée en mémoire pour ne pas les redemander à chaque fois. |
| fetch | La fonction JavaScript qui envoie une requête HTTP depuis le navigateur. |
| XSS | Une attaque où du code (souvent `<script>`) est injecté dans une page. html/template et textContent nous en protègent. |
| Géocodage | Transformer un nom de lieu en coordonnées (latitude, longitude). |
| Leaflet | Une bibliothèque JavaScript pour afficher des cartes interactives. |
