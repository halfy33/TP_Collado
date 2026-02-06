# 📊 Inventory & Buildup - Monitoring Système Multi-Plateforme

Solution complète de monitoring système permettant de collecter et visualiser les métriques système (CPU, mémoire, disques, réseau, processus) de manière centralisée.

## 🎯 Vue d'ensemble

Le projet se compose de deux applications complémentaires :

- **Inventory** : Agent de collecte exposant les métriques système via API REST et interface web
- **Buildup** : Serveur centralisé collectant les données de plusieurs agents Inventory et les stockant dans InfluxDB

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Inventory  │────▶│   Buildup   │────▶│  InfluxDB   │
│   (Agent)   │ GET │ (Serveur)   │ POST│  (Stockage) │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                                              ▼
                                        ┌─────────────┐
                                        │   Grafana   │
                                        │(Visualisation)│
                                        └─────────────┘
```

---

## 📦 Inventory - Agent de Collecte

### Description

Inventory est un agent léger qui collecte les métriques système en temps réel et les expose via :
- **API REST** : Format JSON pour l'intégration avec d'autres systèmes
- **Interface Web** : Tableaux de bord HTML interactifs

### 🚀 Installation

#### Prérequis
- Docker (recommandé) ou Go 1.25+
- Port 80 disponible

#### Déploiement avec Docker

```bash
cd inventory

# Build et lancement
make restart

# Ou manuellement
docker build -t inventory:latest .
docker run -d -p 80:80 --name inventory inventory:latest
```

#### Déploiement sans Docker

```bash
cd inventory
make build
sudo ./bin/inventory
```

### 📚 Utilisation

#### Interface Web

Accédez aux tableaux de bord via votre navigateur :

- **Page principale** : `http://localhost/`
- **Processus** : `http://localhost/procs.html`
  - Visualisation des processus en cours
  - Filtrage par nom, utilisateur, statut
  - **Kill de processus** : bouton 🗑️ pour arrêter un processus
- **Charge système** : `http://localhost/load.html`
  - Load average (1, 5, 15 min)
  - Utilisation CPU par cœur
- **Disques** : `http://localhost/disk.html`
  - Partitions et utilisation
- **Mémoire** : `http://localhost/mem.html`
- **Réseau** : `http://localhost/network.html`

#### API REST

Toutes les métriques sont accessibles en JSON :

| Endpoint | Méthode | Description |
|----------|---------|-------------|
| `/health` | GET | Vérification de l'état du service |
| `/cpu` | GET | Utilisation CPU par cœur |
| `/ps` | GET | Liste de tous les processus |
| `/ps/{user}` | GET | Processus filtrés par utilisateur |
| `/ps/kill/{pid}` | POST | Arrêter un processus |
| `/mem` | GET | Utilisation de la mémoire |
| `/disk` | GET | Utilisation des disques |
| `/avg` | GET | Load average système |
| `/net` | GET | Statistiques réseau globales |
| `/net/{card}` | GET | Statistiques d'une interface réseau |

**Exemples d'utilisation :**

```bash
# Vérifier le service
curl http://localhost/health

# Lister les processus
curl http://localhost/ps

# Voir l'utilisation CPU
curl http://localhost/cpu

# Tuer un processus (PID 1234)
curl -X POST http://localhost/ps/kill/1234

# Processus de l'utilisateur 'root'
curl http://localhost/ps/root
```

### 🔧 Configuration

#### Makefile

```makefile
run       # Exécuter en mode développement
build     # Compiler le binaire
clean     # Supprimer le binaire
image     # Construire l'image Docker
start     # Démarrer le conteneur
stop      # Arrêter le conteneur
restart   # Redémarrer (stop + build + image + start)
```

#### Structure du projet

```
inventory/
├── Dockerfile              # Image Docker
├── Makefile               # Commandes de build
├── bin/                   # Binaires compilés
└── src/
    ├── main.go            # Point d'entrée
    ├── routes.go          # Définition des routes
    ├── handle.go          # Handlers HTTP
    ├── goroutine.go       # Collecte en arrière-plan
    ├── cpu/               # Module CPU
    ├── disk/              # Module disques
    ├── load/              # Module charge système
    ├── memory/            # Module mémoire
    ├── netcard/           # Module réseau
    ├── proc/              # Module processus
    └── www/               # Interface web (HTML/CSS/JS)
```

---

## 🏗️ Buildup - Serveur de Centralisation

### Description

Buildup collecte les métriques de plusieurs agents Inventory et les stocke dans InfluxDB pour analyse et visualisation avec Grafana.

### 🚀 Installation

#### Prérequis
- Docker & Docker Compose
- InfluxDB 2.x
- Ports 8084 (buildup) et 8086 (InfluxDB) disponibles

#### Configuration des serveurs

Éditez `buildup/src/servers.yaml` pour définir les agents à surveiller :

```yaml
servers:
  - "192.168.1.10:80"      # Agent 1
  - "192.168.1.11:80"      # Agent 2
  - "192.168.1.12:80"      # Agent 3
```

#### Variables d'environnement

Créez un fichier `.env` dans `buildup/` :

```env
INFLUXDB_URL=http://influxdb:8086
INFLUXDB_TOKEN=votre-token-influxdb
INFLUXDB_ORG=myorg
INFLUXDB_BUCKET=metrics
TICK_INTERVAL=5
```

#### Déploiement

```bash
cd buildup

# Avec Docker Compose (recommandé)
cd docker
docker-compose up -d

# Ou manuellement
make restart
```

### 📊 Configuration InfluxDB

```bash
# Accéder à l'interface InfluxDB
http://localhost:8086

# Première connexion : créer
# - Organisation : myorg
# - Bucket : metrics
# - Token : Générer et copier dans .env
```

### 📈 Visualisation avec Grafana

```bash
# Accéder à Grafana
http://localhost:3000
# Identifiants par défaut : admin/admin

# Ajouter la source de données InfluxDB
# - Type : InfluxDB
# - URL : http://influxdb:8086
# - Token : votre-token
# - Organisation : myorg
# - Bucket : metrics
```

### 🔄 Architecture de collecte

Buildup interroge les agents Inventory toutes les 5 secondes (configurable) :

1. **Collecte active** : Buildup fait des requêtes GET vers les agents
2. **Stockage** : Données envoyées à InfluxDB
3. **Visualisation** : Grafana affiche les métriques

```go
// Goroutines de collecte par métrique
- GoCPU()    → /cpu
- GoMem()    → /mem
- GoDisk()   → /disk
- GoLoad()   → /avg
- GoNet()    → /net
- GoProcs()  → /ps
```

---

## 🔒 Sécurité

### Recommandations

1. **Firewall** : Limitez l'accès aux ports (80, 8084, 8086)
2. **Authentification** : Ajoutez un reverse proxy (nginx) avec authentification
3. **HTTPS** : Utilisez des certificats SSL/TLS en production
4. **Permissions** : Lancez inventory avec les permissions minimales nécessaires

### Exemple avec nginx (reverse proxy)

```nginx
server {
    listen 443 ssl;
    server_name inventory.example.com;

    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/private/key.pem;

    location / {
        proxy_pass http://localhost:80;
        auth_basic "Restricted";
        auth_basic_user_file /etc/nginx/.htpasswd;
    }
}
```

---

## 🌍 Déploiement Multi-Plateforme

### Linux (déjà déployé)
```bash
make build
```

### Windows
```bash
GOOS=windows GOARCH=amd64 make build
# Généré : bin/inventory.exe
```

### macOS
```bash
GOOS=darwin GOARCH=amd64 make build
# Généré : bin/inventory (macOS)
```

### ARM (Raspberry Pi)
```bash
GOOS=linux GOARCH=arm64 make build
```

---

## 🐛 Dépannage

### Inventory ne démarre pas

```bash
# Vérifier les logs
docker logs inventory

# Vérifier le port
sudo netstat -tlnp | grep :80

# Rebuild complet
make clean && make restart
```

### Erreurs de compilation

```bash
# Synchroniser les dépendances
cd src
go mod tidy
go mod download
```

### Buildup n'atteint pas les agents

```bash
# Tester la connectivité
curl http://agent-ip:80/health

# Vérifier servers.yaml
cat buildup/src/servers.yaml

# Vérifier les logs
docker logs buildup
```

### Données manquantes dans InfluxDB

```bash
# Vérifier la connexion
docker exec buildup sh -c 'ping influxdb'

# Tester le token
curl -H "Authorization: Token VOTRE_TOKEN" \
  http://localhost:8086/api/v2/buckets
```

---

## 📝 Notes importantes

### Kill de processus

- ⚠️ **Attention** : Tuer des processus système peut rendre le système instable
- 🔒 Nécessite les permissions appropriées
- 🐳 Dans Docker : seuls les processus du conteneur sont visibles
- 💡 Conseil : Tester avec des processus non critiques (`sleep`, etc.)

### Performance

- **Intervalle de collecte** : 5 secondes par défaut
- **Impact système** : Minimal (~0.5% CPU, 10-20 MB RAM)
- **Bande passante** : ~2-5 KB/s par métrique

### Limitations connues

- Dans Docker : isolation des processus (seuls ceux du conteneur sont visibles)
- Processus zombies : ne peuvent pas être tués
- PID 1 : processus principal, ne peut être tué

---


## Auteurs
* **Lucas** - [Luas-IQ21](https://github.com/Lucas-IQ21)
* **Yanis** - [halfy33](https://github.com/halfy33)
* **Thomas** - [ThomasCelle](https://github.com/ThomasCelle)
* **Romain** - [RomainBnr](https://github.com/RomainBnr)


## 🎓 Ressources

- [Documentation Go](https://golang.org/doc/)
- [gopsutil](https://github.com/shirou/gopsutil)
- [InfluxDB](https://docs.influxdata.com/influxdb/v2/)
- [Grafana](https://grafana.com/docs/)
- [Docker](https://docs.docker.com/)

---

**Version** : 1.0  
**Date** : Février 2026
