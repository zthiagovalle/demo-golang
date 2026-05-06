# demo-golang — Clean Architecture / DDD em Go

Demo didática de uma arquitetura em camadas (Clean Architecture + DDD) organizada **por feature**, em Go com Echo v5, Postgres e SNS/SQS via LocalStack. Os princípios se aplicam a qualquer linguagem (Java/Spring, TS/Nest, Python/FastAPI, C#, Kotlin, etc.).

---

## 1. A regra de ouro

> Dependências apontam **sempre para o domínio**. `core` é transversal — todas as camadas importam dele, ele não importa de ninguém.

```
            ┌──────►  core  ◄──────┐
            │     (transversal)    │
            │                      │
application  ──►   domain   ◄──  infra
 (entrada)        (núcleo)        (saída)
```

- **`domain`** não importa nada de fora (só `core`). Define os **contratos** (`I{Verb}{Entity}Usecase`, `I{Entity}Repository`, `I{Topic}Producer`) que precisa.
- **`infra`** implementa esses contratos (Postgres, SNS, HTTP externo, ...).
- **`application`** orquestra: recebe input (HTTP, SQS), chama use case, devolve output.
- **`core`** contém utilidades transversais (config, validator wrapper, route helper, response helpers). **Regra inversa**: `core` não pode importar de `domain`, `application` ou `infra` — só de stdlib e libs externas. Isso garante que ele permaneça reutilizável e neutro de domínio.

Isso é Dependency Inversion. Trocar Postgres por DynamoDB, ou SNS por Kafka, é uma mudança contida em `infra/` — `domain/` não muda uma linha.

---

## 2. Tour pelas camadas

### `src/domain/` — núcleo de negócio
| Pasta | O que tem | Exemplo |
|---|---|---|
| `models/` | Entidades + DTOs (`_create`, `_update`, `_page_params`) com tags `validate:"..."` | `product.go`, `product_create.go` |
| `enums/`  | Tipos enumerados com `String()` + função validadora p/ `validator/v10` | `product_status.go` (ACTIVE/INACTIVE + `ProductStatusValidator`) |
| `exceptions/` | Códigos de erro como `const string` | `ErrProductNotFound = "errProductNotFound"` |
| `repositories/` | **Ports** — interfaces de persistência (`I{Entity}Repository`) | `products_repository.go` (só `IProductsRepository`) |
| `producers/`    | **Ports** — interfaces de publicação de eventos (`I{Topic}Producer`) | `product_producer.go` (só `IProductProducer`) |
| `usecases/` | **Um arquivo por caso de uso** com a interface `I{Verb}{Entity}Usecase` declarada no próprio arquivo + diretiva `//go:generate mockgen` | `create_product_usecase.go` |
| `*/mock/` | Mocks **gerados** via `mockgen` (pacotes `usecasesmock`, `repositoriesmock`, `producersmock`) | `create_product_usecase_mock.go` |

> **Ports vivem em `domain/`** (Hexagonal puro). Use cases dependem só de `domain/...` — zero import de `infra/`. Não há `ports.go` único: cada interface fica no seu próprio arquivo, próximo das que vivem na mesma camada conceitual.

### `src/application/` — adaptadores de entrada
| Pasta | Função |
|---|---|
| `controllers/` | HTTP via Echo. Faz **bind → validate → use case → response**. Expõe `Routes() []shared.Route` que o `main.go` itera para registrar. Anotações Swagger nos comentários de cada handler. |
| `consumers/` | Consome eventos da fila SQS e dispara orquestração. |

### `src/infra/` — adaptadores de saída
| Pasta | Adapter | Notas |
|---|---|---|
| `repositories/products_postgres_repository.go` | `ProductsPostgresRepository` implementa `domain/repositories.IProductsRepository` | SQL como `const` no topo do arquivo (`insertProductQuery`, etc.) |
| `producers/product_sns_producer.go` | `ProductSNSProducer` implementa `domain/producers.IProductProducer` | publica JSON no SNS |
| `gateways/` | Cliente HTTP de serviço externo (catálogo) | mockado via WireMock |

> O nome do adapter explicita a tecnologia (`...Postgres`, `...SNS`). Trocar Postgres por DynamoDB = adicionar `ProductsDynamoRepository` ao lado e trocar 1 linha no `main.go` — sem mexer em domain/usecases.

### `src/core/` — utilidades transversais
- `config/` — carrega env vars uma vez no boot.
- `shared/` — `Route`/`RegisterRoutes` (helper p/ Echo), `CustomValidator` (wrapper de `validator/v10`), `ErrorJSON` / `InternalErrorJSON`.

> **Direção da dependência**: `domain`, `application` e `infra` podem importar de `core`. `core` **nunca** importa dessas camadas — depende só de stdlib e libs externas. Se você se pegar tentando referenciar um model ou exception aqui, é sinal de que o helper deveria viver em outro lugar.

### `main.go` — composition root
**Único** lugar que conhece as implementações concretas. Faz o wiring (DI explícita), registra validators custom (`oneOfProductStatus`) e registra rotas via `shared.RegisterRoutes(e, productsController.Routes())`.

---

## 3. Convenções de código (estilo da equipe Zaelotech)

| Convenção | Aplicação |
|---|---|
| Interfaces com prefixo `I` | `ICreateProductUsecase`, `IProductsRepository`, `IProductProducer` |
| Adapter nomeado pela tecnologia | `ProductsPostgresRepository`, `ProductSNSProducer` (deixa óbvio o backend) |
| Pluralização da entidade | `products_repository.go`, `products_v1_controller.go` |
| Erros como `const string` (não tipados) | `errors.New(exceptions.ErrProductNotFound)` no use case; controller compara `err.Error() == exceptions.ErrXyz` para mapear status HTTP |
| SQL como `const` no topo do repositório | `insertProductQuery`, `findAllPaginatedProductsQuery`, ... |
| `Routes()` declarativo | Controller devolve slice de rotas; main.go registra |
| Anotações Swagger em comentário | `@Summary`, `@Tags`, `@Router` em cada handler — `make doc` gera o JSON |
| Validação na entrada | tags `validate:"..."` nos DTOs + `c.Validate(&input)` no controller |

---

## 4. Por que organizar por feature, não por tipo

Estrutura clássica "por tipo": `controllers/ services/ repositories/ models/` — para mexer em "produto" você abre 4 pastas.

Estrutura por feature (essa demo):
```
domain/usecases/create_product_usecase.go
domain/usecases/update_product_usecase.go
domain/repositories/products_repository.go         (interface)
infra/repositories/products_postgres_repository.go (adapter)
application/controllers/products_v1_controller.go
```
**Tudo de produto vive próximo.** Onboarding mais rápido, menos merge conflicts.

---

## 5. Endpoints expostos

| Método | Rota | Use case |
|---|---|---|
| GET    | `/v1/products?page=1&pageSize=20` | get_all_paginated_products |
| POST   | `/v1/products` | create_product (publica `product.created` no SNS) |
| PUT    | `/v1/products/:id` | update_product |
| DELETE | `/v1/products/:id` | delete_product |
| PATCH  | `/v1/products/:id/toggle-status` | toggle_product_status |

Em paralelo, um **consumer** lê `product-created-queue` (subscrita ao tópico SNS) e chama o **catalog-gateway** (WireMock) para enriquecimento — exercitando o caminho assíncrono completo.

---

## 6. Como rodar

Pré-requisitos: Docker, Go 1.26+.

```bash
# 1. (uma vez) instalar tooling: mockgen, go-acc, swag
make install-dependencies

# 2. (uma vez por clone) gera mocks + tidy + cria .env
make init

# 3. subir Postgres + pgAdmin + LocalStack + WireMock
make env-up
# Abre o pgAdmin em http://localhost:5050 — server "demo-postgres" já pré-configurado
# (sem login: PGADMIN_CONFIG_SERVER_MODE=False)

# 4. aplicar migrations
make migration-up

# 5. (opcional) seed de dados
make seed

# 6. subir o app
make run

# 7. testar (ver coleção Bruno em /Users/thiago/www/bruno/demo-golang)
curl -X POST http://localhost:8080/v1/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"Cafe Especial","description":"Grao premium","price_cents":4990}'

curl 'http://localhost:8080/v1/products?page=1&pageSize=20'

# 8. ver mensagem na fila
cd development-environment && make awslocal-receive

# 9. rodar testes (regenera mocks + go-acc + filtra main/_mock)
#    Os testes de repositório sobem um Postgres efêmero via testcontainers — Docker precisa estar rodando.
make test
```

Outros targets úteis: `make mock` (regenera mocks), `make doc` (Swagger JSON), `make deps` (atualiza tudo), `make fmt`, `make cover` (HTML coverage), `make build*` (binários multi-plataforma).

---

## 7. Mocks com `mockgen`

Toda interface tem uma diretiva no topo do arquivo:

```go
//go:generate mockgen -source create_product_usecase.go -destination mock/create_product_usecase_mock.go -package usecasesmock
```

`make mock` apaga todos os `*_mock.go` e roda `go generate ./...`. O resultado:

```
src/domain/usecases/mock/        → pacote usecasesmock
src/domain/repositories/mock/    → pacote repositoriesmock
src/domain/producers/mock/       → pacote producersmock
```

Uso típico em teste:
```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()

repo := repositoriesmock.NewMockIProductsRepository(ctrl)
repo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

uc := usecases.NewCreateProductUsecase(repo, producer)
```

> Nota: a equipe Zaelotech usa `github.com/golang/mock` v1.6.0 (deprecated). Aqui adotamos o fork mantido `go.uber.org/mock` (API idêntica, só muda o import path do pacote `gomock`).

---

## 8. Tradeoffs

**A favor:**
- Testabilidade altíssima — domain testa sem banco, sem rede, em milissegundos.
- Trocar infra é barato e localizado.
- Múltiplos times trabalham por feature sem pisar uns nos outros.
- Use cases pequenos e focados são fáceis de revisar e dar ownership.

**Contra:**
- Boilerplate: cada dependência vira interface + impl + mock gerado.
- Curva inicial — devs novos precisam internalizar a regra de dependência e o pattern de mocks.
- **Overkill** para CRUDs triviais, MVPs descartáveis ou scripts.

**Quando aplicar:** serviço com regras de negócio reais, vida longa, múltiplos contribuidores. **Quando não:** prototipo de 200 linhas, ETL pontual, lambda simples.

---

## 9. Replicando em outras linguagens

| Conceito Go | Java/Spring | TS/NestJS | Python/FastAPI |
|---|---|---|---|
| Interface `I{Verb}{Entity}Usecase` no próprio arquivo | interface + `@Component` | abstract class + DI token | Protocol / ABC |
| Struct + construtor com DI | `@Service` + `@Autowired` | `@Injectable()` + constructor | classe + injeção manual ou via container |
| `main.go` wiring | Spring Boot (auto-config) | módulos NestJS | factory function no startup |
| `mockgen` + `go generate` | Mockito | jest.mock / ts-mockito | unittest.mock |
| `validator/v10` tags | `jakarta.validation` | `class-validator` | `pydantic` |

A **regra de dependência e a separação por camada/feature** são iguais.

---

## 10. O que **não** está nesta demo (de propósito)

- Auth / middleware de segurança
- Observabilidade (métricas, tracing, logs estruturados completos)
- CI/CD
- Suite completa de testes (cobertura focada — temos use case, controller e repositório com testcontainers; faltam consumer, producer, gateway)
- Geração efetiva do Swagger (anotações estão prontas; basta `make doc` após instalar `swag`)

Foram deixados de fora para manter o foco arquitetural. Em produção, todos seriam adicionados sem mudar a estrutura.
