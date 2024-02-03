package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func InitRouter() (*gin.Engine) {
    
    router := gin.Default();
    router.Use(cors.Default());

    // Bootstrap Endpoints of API
    // ==========================
    router.POST("/blockchat_api/register_node", RegisterNode);

    // Backend Endpoints of API
    // ========================
    router.POST("/blockchat_api/receive_new_node", ReceiveNewNode);
    router.POST("/blockchat_api/validate_transaction", ValidateTransaction);
    router.POST("/blockchat_api/receive_transaction", ReceiveTransaction);
    router.POST("/blockchat_api/receive_stake", ReceiveStake);
    router.POST("/blockchat_api/validate_stake", ValidateStake);

    // Client Endpoints of API
    // =======================
    router.POST("/blockchat_api/set_stake", SetStake);
    router.POST("/blockchat_api/send_transaction", SendTransaction);
    router.GET("/blockchat_api/get_balance", GetBalance);
    router.GET("/blockchat_api/get_last_block", GetLastBlock);

    return router;
}


