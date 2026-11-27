import os
import traceback

import pymongo
from dotenv import load_dotenv
from flask import Flask, request, jsonify, g
from flask_cors import CORS

from core.api.health_check_api import HealthCheckAPI
from web import Web

app = Flask(__name__)
CORS(app)
load_dotenv()

jwt_signing_key = os.getenv("JWT_SIGNING_KEY", "secret")

db = pymongo.MongoClient(
    os.getenv("DB_HOST", "localhost"),
    int(os.getenv("DB_PORT", "27017"))
)[os.getenv("DB_NAME", "vertex")]

Web(app).register()

HealthCheckAPI(app).register()


@app.before_request
def before_request():
    g.request_body = request.get_json(force=True, silent=True)
    g.jwt_signing_key = jwt_signing_key


@app.errorhandler(404)
def handle_404_error(e):
    return jsonify({
        "success": False,
        "message": str(e)
    }), 404


@app.errorhandler(Exception)
def handle_all_errors(e):
    print(traceback.format_exc())
    return jsonify({
        "success": False,
        "message": str(e)
    }), 500


if __name__ == '__main__':
    app.run(
        host=os.getenv("HOST", "0.0.0.0"),
        port=int(os.getenv("PORT", "5000")),
        debug=os.getenv("ENV", "DEV") != "PROD"
    )
