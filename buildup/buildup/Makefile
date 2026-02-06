# Variables
COMPOSE         ?= docker compose      
DOCKER_DIR      ?= docker
COMPOSE_FILE    ?= docker-compose.yml

run:
	@cd src ; go run main.go

build:
	@cd src ; go mod tidy ; go build -o ../bin/buildup *.go

clean:
	@rm ./bin/buildup

image:
	@docker build -t buildup:latest .

start:
	@docker run -d -p 80:80 --name buildup buildup:latest

stop:
	@docker rm -f buildup

# Démarrage stack
up:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) up -d

# Arrêt de la stack
down:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) down

# Journaux de la stack
logs:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) logs -f

# Liste des processus de la stack
ps:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) ps