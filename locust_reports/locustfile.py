from locust import HttpUser, task

class HelloWorldUser(HttpUser):
    @task
    def biz(self):
        self.client.post(url="/biz/get_users", json={
    		"user_id": 1,
    		"auth_key": 6,
    		"message_id":6
	})
        self.client.post(url="/biz/get_users_with_sql_inject", json={
                "user_id": "1",
                "auth_key": 6,
                "message_id":6
        })
    @task
    def auth(self):
        self.client.post(url="/auth/req_pq", json={
	 	"nonce":"fdvefverfbefb",
    		"message_id":2
        })
        self.client.post(url="/auth/req_DH_pq", json={
    		"nonce":"fdvefverfbefb",
    		"server_nonce": "ORfAZdBzINgDmSrAgbXp",
    		"message_id":4,
    		"b":5
        })
