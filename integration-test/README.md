Сборка образов
Docker build -t user-cart-order-image:latest -f user-cart-order/Dockerfile .

Docker build -t product-image:latest -f product/Dockerfile .

docker build -t llm-service-image:latest -f llm-service/Dockerfile llm-service/
