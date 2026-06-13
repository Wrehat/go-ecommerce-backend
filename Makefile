# Ganti <NAMA_POSTGRES> dan <NAMA_REDIS> dengan nama asli di Docker Desktop-mu
infra-up:
	docker start postgresKana
	docker start redis-ecommerce

infra-down:
	docker stop postgresKana
	docker stop redis-ecommerce

# Perintah sakti untuk menjalankan semuanya
run: infra-up
	go run cmd/api/main.go