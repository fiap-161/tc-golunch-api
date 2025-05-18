echo "Derrubando containers, volumes e redes..."
docker-compose down -v --remove-orphans

echo "Reconstruindo e subindo os containers..."
docker-compose up --build