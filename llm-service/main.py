from __future__ import annotations
from fastapi import FastAPI, Request, HTTPException
import os
import redis
import json
from dotenv import load_dotenv
from yandex_cloud_ml_sdk import YCloudML
from typing import Optional

# Загружаем переменные окружения из файла .env
load_dotenv()

# Получение конфигурационных переменных
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
YANDEX_API_KEY = os.getenv('YANDEX_API_KEY')
YANDEX_FOLDER_ID = os.getenv('YANDEX_FOLDER_ID')
MODEL_NAME = os.getenv('MODEL_NAME', 'yandexgpt')

# Проверка обязательных переменных
if not all([YANDEX_API_KEY, YANDEX_FOLDER_ID]):
    raise RuntimeError("YANDEX_API_KEY and YANDEX_FOLDER_ID must be set in .env")

# Инициализация SDK Yandex Cloud ML
yandex_sdk = YCloudML(
    folder_id=YANDEX_FOLDER_ID,
    auth=YANDEX_API_KEY
)

# Подключение к Redis
redis_client = redis.Redis(
    host=REDIS_HOST,
    port=REDIS_PORT,
    db=0
)

app = FastAPI()

def get_dialog_history(user_id: str, product_id: Optional[str] = None) -> list:
    """Получает историю диалога из Redis"""
    key = f"{user_id}:{product_id}" if product_id else user_id
    history = redis_client.get(key)
    return json.loads(history) if history else []

def save_dialog_history(user_id: str, product_id: Optional[str], history: list):
    """Сохраняет историю диалога в Redis"""
    key = f"{user_id}:{product_id}" if product_id else user_id
    redis_client.set(key, json.dumps(history), ex=86400)  # TTL 24 часа

@app.post("/chat")
async def chat_endpoint(request: Request):
    """Обработчик чат-запросов"""
    try:
        data = await request.json()

        # Проверка обязательных полей
        if 'user_id' not in data or 'question' not in data:
            raise HTTPException(status_code=400, detail="Требуются user_id и question")

        user_id = data['user_id']
        question = data['question']
        product_id = data.get('product_id')

        # Получение истории диалога
        history = get_dialog_history(user_id, product_id)

        # Формируем сообщения для модели
        messages = [
            {
                "role": "system",
                "text": "Ты — консультант магазина."
            }
        ]

        # Добавляем историю диалога
        for i, msg in enumerate(history):
            role = "user" if i % 2 == 0 else "assistant"
            messages.append({"role": role, "text": msg})

        # Добавляем текущий вопрос
        messages.append({
            "role": "user",
            "text": question
        })

        # Запрос к Yandex GPT
        result = (
            yandex_sdk.models
            .completions(MODEL_NAME)
            .configure(temperature=0.6)
            .run(messages)
        )

        # Получение ответа
        response_text = result[0].text

        # Обновление истории диалога
        save_dialog_history(
            user_id,
            product_id,
            history + [question, response_text]
        )

        return {"answer": response_text}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    """Проверка состояния сервиса"""
    try:
        # Проверка соединения с Redis
        if not redis_client.ping():
            raise Exception("Redis connection failed")
        # Проверка работы Yandex GPT
        test_messages = [{"role": "user", "text": "Тестовое сообщение"}]
        yandex_sdk.models.completions(MODEL_NAME).run(test_messages)
        return {"status": "OK"}
    except Exception as e:
        raise HTTPException(status_code=503, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
