from fastapi import FastAPI
from handler import router
import sys
import os

app = FastAPI()
app.include_router(router)

if __name__ == "__main__":
    import uvicorn
    sys.path.insert(0, os.path.join(os.path.dirname(__file__), "proto"))
    uvicorn.run(app, host="0.0.0.0", port=8000)
