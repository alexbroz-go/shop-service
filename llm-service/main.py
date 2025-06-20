from __future__ import annotations
from fastapi import FastAPI, Request, HTTPException
from datetime import datetime
import os
import redis
import json
from dotenv import load_dotenv
from yandex_cloud_ml_sdk import YCloudML
from typing import Optional, List, Dict
from pydantic import BaseModel, validator

# Загрузка переменных окружения
load_dotenv()

# Модели данных
class ChatMessage(BaseModel):
    role: str  # "system", "user", "assistant"
    text: str
    timestamp: Optional[str] = None
    tool_used: Optional[str] = None
    sources: Optional[List[str]] = None

    @validator('role')
    def validate_role(cls, v):
        if v not in {"system", "user", "assistant"}:
            raise ValueError("Invalid message role")
        return v

# Конфигурация
REDIS_HOST = os.getenv('REDIS_HOST', 'redis_db')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
YANDEX_API_KEY = os.getenv('YANDEX_API_KEY')
YANDEX_FOLDER_ID = os.getenv('YANDEX_FOLDER_ID')
MODEL_NAME = os.getenv('MODEL_NAME', 'yandexgpt')

# Проверка конфигурации
if not all([YANDEX_API_KEY, YANDEX_FOLDER_ID]):
    raise RuntimeError("Missing required environment variables")

# Инициализация клиентов
yandex_sdk = YCloudML(
    folder_id=YANDEX_FOLDER_ID,
    auth=YANDEX_API_KEY
)

redis_client = redis.Redis(
    host=REDIS_HOST,
    port=REDIS_PORT,
    db=0,
    decode_responses=False
)

app = FastAPI()

# Системный промпт для строительного магазина
SYSTEM_PROMPT = {
    "role": "system",
    "text": """Ты — консультант магазина строительных инструментов. Правила:
1. Давай точные технические характеристики
2. Уточняй тип работ (проф/бытовые) 
3. Рекомендуй только реальные товары
4. Предупреждай о необходимости СИЗ
5. Не давай советов по монтажу
6. Не давай совет по моделям инструмента, которые нужно покупать клиенту, если он сам не спросил про модель"""
}

def get_dialog_history(user_id: str, product_id: Optional[str] = None) -> List[Dict]:
    """Получает полную историю диалога с сохранением структуры сообщений"""
    key = f"chat:{user_id}:{product_id}" if product_id else f"chat:{user_id}"
    try:
        history = redis_client.get(key)
        return json.loads(history) if history else []
    except json.JSONDecodeError:
        return []

def save_dialog_history(
    user_id: str,
    messages: List[Dict],
    product_id: Optional[str] = None,
    ttl: int = 86400
):
    """Сохраняет историю диалога с метаданными"""
    key = f"chat:{user_id}:{product_id}" if product_id else f"chat:{user_id}"
    try:
        redis_client.set(key, json.dumps(messages), ex=ttl)
    except Exception as e:
        print(f"Error saving history: {e}")

@app.post("/chat")
async def chat_endpoint(request: Request):
    try:
        data = await request.json()
        
        if 'user_id' not in data or 'question' not in data:
            raise HTTPException(status_code=400, detail="Missing required fields")
        
        user_id = data['user_id']
        question = data['question']
        product_id = data.get('product_id')

        history = get_dialog_history(user_id, product_id)
        messages = [SYSTEM_PROMPT] + history
        
        new_message = {
            "role": "user",
            "text": question,
            "timestamp": datetime.utcnow().isoformat()
        }
        messages.append(new_message)
        
        messages_for_api = [
            {"role": m["role"], "text": m["text"]} 
            for m in messages
            if "role" in m and "text" in m
        ]
        
        result = (
            yandex_sdk.models
            .completions(MODEL_NAME)
            .configure(temperature=0.6)
            .run(messages_for_api)
        )
        
        if not result:
            raise HTTPException(status_code=500, detail="Empty model response")

        # Очистка ответа от переносов строк
        response_text = result[0].text.replace('\n', ' ')
        
        assistant_message = {
            "role": "assistant",
            "text": response_text,
            "timestamp": datetime.utcnow().isoformat()
        }
        updated_history = history + [new_message, assistant_message]
        save_dialog_history(user_id, updated_history, product_id)
        
        return {"answer": response_text}

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    """Проверка работоспособности сервиса"""
    try:
        if not redis_client.ping():
            raise Exception("Redis unavailable")
        
        test_messages = [{"role": "user", "text": "Тест соединения"}]
        yandex_sdk.models.completions(MODEL_NAME).run(test_messages)
        
        return {
            "status": "OK",
            "redis": True,
            "yandex_gpt": True
        }
    except Exception as e:
        raise HTTPException(status_code=503, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)