from fastapi import APIRouter, Request, HTTPException
from datetime import datetime
from models import ChatMessage, SYSTEM_PROMPT
from dialog_utils import get_dialog_history, save_dialog_history
from gpt import get_gpt_response
from product import get_product_info

router = APIRouter()

@router.post("/chat")
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

        if product_id:
            product_info = get_product_info(product_id)
            if product_info:
                messages.append({
                    "role": "system",
                    "text": f"Информация о товаре: {product_info}"
                })

        messages_for_api = [
            {"role": m["role"], "text": m["text"]}
            for m in messages if "role" in m and "text" in m
        ]

        response_text = get_gpt_response(messages_for_api)

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
