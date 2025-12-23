git pull origin main
go mod tidy
GOOS=linux GOARCH=amd64 go build -o checkip
pm2 restart checkip