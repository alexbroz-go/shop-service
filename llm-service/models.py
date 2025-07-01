from typing import Optional, List
from pydantic import BaseModel, validator

class ChatMessage(BaseModel):
    role: str
    text: str
    timestamp: Optional[str] = None
    tool_used: Optional[str] = None
    sources: Optional[List[str]] = None

    @validator('role')
    def validate_role(cls, v):
        if v not in {"system", "user", "assistant"}:
            raise ValueError("Invalid message role")
        return v

SYSTEM_PROMPT = {
    "role": "system",
    "text": (
        "Ты — консультант магазина строительных инструментов. Правила:\n"
        "1. Давай точные технические характеристики\n"
        "2. Уточняй тип работ (проф/бытовые)\n"
        "3. Рекомендуй только реальные товары\n"
        "4. Предупреждай о необходимости СИЗ\n"
        "5. Не давай советов по монтажу\n"
        "6. Не давай совет по моделям инструмента, которые нужно покупать клиенту, если он сам не спросил про модель"
    )
}
