rem docker kill $(docker ps -q)
rem docker rm $(docker ps -a -q)
docker rm imagemaker-app-run
docker rmi imagemaker-app
docker build -t imagemaker-app .
docker run --publish 8082:8000 --rm --name imagemaker-app-run imagemaker-app
explorer "http://localhost:8082"