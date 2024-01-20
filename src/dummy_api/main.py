from flask_restful import Api

from dummy_api import app
from .auth import *
api = Api(app, prefix='/blockchat_api')

api.add_resource(SendTransaction, '/send_transaction')
api.add_resource(SetStake, '/set_stake')
api.add_resource(GetLastBlock, '/get_last_block')
api.add_resource(GetBalance, '/get_balance')


if __name__ == "__main__":  
    # debug mode: every time we change it restarts the server
    app.run(debug=True, port=9876)
