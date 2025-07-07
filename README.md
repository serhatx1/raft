# Raft Consensus Algorithm (Go)

Bu proje, Go ile yazılmış basit bir Raft consensus algoritması örneğidir.

## Başlatma (Docker ile)

Projeyi Docker ve docker-compose ile kolayca başlatabilirsiniz.

### 1. Docker image oluşturun ve node'ları başlatın

```sh
docker-compose up --build
```

Bu komut, üç adet raft node'u başlatır:
- Node 1: http://localhost:8001
- Node 2: http://localhost:8002
- Node 3: http://localhost:8003

### 2. API Uç Noktaları

Her node aşağıdaki HTTP endpoint'lerini sunar:
- `/` : Basit karşılama mesajı
- `/request-vote` : Raft seçim isteği (POST)
- `/heartbeat` : Raft heartbeat (POST)

### 3. Ortam Değişkenleri

docker-compose ile başlatırken her node için şu ortam değişkenleri ayarlanır:
- `NODE_ID` : Node'un benzersiz ID'si
- `NODE_ADDR` : Node'un dinlediği adres (örn. :8001)
- `NODE_PEERS` : Virgülle ayrılmış diğer node adresleri (örn. http://node2:8002,http://node3:8003)

> Not: Varsayılan olarak, kodun environment değişkenlerini okuması için main.go'da güncelleme yapmanız gerekebilir.

## Manuel Başlatma (Geliştiriciler için)

Go yüklü ise, projeyi manuel başlatmak için:

```sh
go run main.go
```

## Lisans

MIT

## Klasör Yapısı

- `cmd/node.go`: Node başlatma, RPC server/client
- `internal/raft/state.go`: Node state ve term yönetimi
- `internal/raft/rpc.go`: RequestVote, AppendEntries RPC
- `internal/raft/election.go`: Lider seçimi logic
- `internal/raft/log.go`: Log entry yönetimi
- `internal/statemachine.go`: Commit edilen komutları işleyen logic
- `test/`: Unit ve integration testler

## Kurulum

```sh
go mod tidy
go build ./cmd/node.go
```

## Açıklama

Her dosya ve klasör, Raft algoritmasının farklı bir sorumluluğunu üstlenir. Detaylı açıklamalar ilgili dosyalarda yer alacaktır.
