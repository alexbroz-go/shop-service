from yandex_cloud_ml_sdk import YCloudML
from config import YANDEX_API_KEY, YANDEX_FOLDER_ID, MODEL_NAME

yandex_sdk = YCloudML(folder_id=YANDEX_FOLDER_ID, auth=YANDEX_API_KEY)

def get_gpt_response(messages):
    result = (
        yandex_sdk.models
        .completions(MODEL_NAME)
        .configure(temperature=0.6)
        .run(messages)
    )
    if not result:
        raise RuntimeError("Empty response from YandexGPT")
    return result[0].text.replace('\n', ' ')
