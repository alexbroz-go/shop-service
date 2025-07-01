import os
from dotenv import load_dotenv

load_dotenv()

REDIS_HOST = os.getenv('REDIS_HOST', 'redis_db')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
YANDEX_API_KEY = os.getenv('YANDEX_API_KEY')
YANDEX_FOLDER_ID = os.getenv('YANDEX_FOLDER_ID')
MODEL_NAME = os.getenv('MODEL_NAME', 'yandexgpt')
GRPC_PRODUCT_HOST = os.getenv('PRODUCT_GRPC_HOST', 'product:50051')
