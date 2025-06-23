FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

COPY . .

RUN python -m grpc_tools.protoc -Iproto --python_out=proto --grpc_python_out=proto proto/product.proto && \
    sed -i 's/import product_pb2 as product__pb2/from . import product_pb2 as product__pb2/' proto/product_pb2_grpc.py

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
