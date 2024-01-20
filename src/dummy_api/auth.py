from flask import request, jsonify, make_response, session
from functools import wraps
from flask_restful import Resource
import json, jwt, datetime
from dummy_api import app
from flask_cors import cross_origin,CORS

CORS(app, resources={r"/blockchat_api/*": {"origins": "http://127.0.0.1:9876"}})

class SendTransaction(Resource):
    @cross_origin()
    def post(self):
        # The id of the node that will receive the message or bitcoin
        recipient_id = request.form.get("recipient_id")
        # Boolean, 0 for message, 1 for bitcoin
        message_or_bitcoin = request.form.get("message_or_bitcoin")
        # Data
        data = request.form.get("data")

        print(recipient_id)
        print(message_or_bitcoin)
        print(data)

        return {'status':'koble'},200

# address is http://127.0.0.1:9876/blockchat_api/set_stake
class SetStake(Resource):
    @cross_origin()
    def post(self):
        # Amount to stake
        stake = request.form.get("stake")

        print(stake)

        return {'status':'koble'},200

class GetBalance(Resource):
    @cross_origin()
    def get(self):
        return {'balance': 25},200

class GetLastBlock(Resource):
    @cross_origin()
    def get(self):
        last_block = {"index": 1234,
                      "previous_hash": "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
                      "transactions": [
                        "tx1hash",
                        "tx2hash",
                        "tx3hash"
                      ],
                      "validator": "validator123",
                      "timestamp": "2024-01-19T12:00:00Z"}

        result = json.dumps(last_block)
        print(str(result))

        return {'last_block':result},200







