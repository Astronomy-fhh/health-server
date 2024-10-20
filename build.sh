docker build -t health-server:latest .

docker save -o health-server.tar health-server:latest

scp health-server.tar ecs:/root/health/run/server


