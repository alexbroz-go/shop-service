import json
from redis_service import redis_client
from typing import Optional, List, Dict

def get_dialog_history(user_id: str, product_id: Optional[str] = None) -> List[Dict]:
    key = f"chat:{user_id}:{product_id}" if product_id else f"chat:{user_id}"
    try:
        history = redis_client.get(key)
        return json.loads(history) if history else []
    except json.JSONDecodeError:
        return []

def save_dialog_history(user_id: str, messages: List[Dict], product_id: Optional[str] = None, ttl: int = 86400):
    key = f"chat:{user_id}:{product_id}" if product_id else f"chat:{user_id}"
    try:
        redis_client.set(key, json.dumps(messages), ex=ttl)
    except Exception as e:
        print(f"Error saving history: {e}")
