import time
import uuid
from functools import wraps

import jwt
from flask import g, request


def authenticate_user(f):
    @wraps(f)
    def decorator(*args, **kwargs):
        token = None

        if 'Authorization' in request.headers:
            auth_header = request.headers['Authorization']
            if auth_header.startswith('Bearer '):
                token = auth_header.split(' ')[1]

        if not token:
            raise Exception("Invalid access token")

        try:
            vertex_endpoint = g.config_service.get_vertex_endpoint()
            jwt_signing_key = g.config_service.get_jwt_signing_key()

            data = jwt.decode(
                token,
                jwt_signing_key,
                algorithms=['HS256'],
                issuer=vertex_endpoint,
                audience=vertex_endpoint
            )

            if data["type"] != "user":
                raise Exception("Invalid access token")

            current_user = {
                "username": data["sub"]
            }

        except Exception as e:
            print(e)
            raise Exception("Invalid access token")

        return f(current_user, *args, **kwargs)

    return decorator



def generate_user_token(username: str) -> str:
    vertex_endpoint = g.config_service.get_vertex_endpoint()
    jwt_signing_key = g.config_service.get_jwt_signing_key()

    payload = {
        "sub": username,
        "type": "user",
        "iss": vertex_endpoint,
        "aud": vertex_endpoint,
        "iat": int(time.time()),
        "exp": int(time.time()) + 7 * 24 * 60 * 60,
        "jti": str(uuid.uuid4())
    }

    token = jwt.encode(payload, jwt_signing_key, algorithm='HS256')
    return token
