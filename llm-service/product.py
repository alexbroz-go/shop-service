# gRPC-клиент для product-сервиса
import grpc
from config import GRPC_PRODUCT_HOST
from proto import product_pb2 as product__pb2
from proto.product_pb2_grpc import ProductServiceStub

channel = grpc.insecure_channel(GRPC_PRODUCT_HOST)
product_stub = ProductServiceStub(channel)

def get_product_info(product_id: str) -> str:
    try:
        product_id_int = int(product_id)
        response = product_stub.GetProduct(product__pb2.GetProductRequest(product_id=product_id_int))
        return f"{response.title} — {response.description}. Цена: {response.price}₽"
    except grpc.RpcError as e:
        print(f"gRPC error: {e}")
        return ""
    except ValueError:
        print("Invalid product_id format")
        return ""
