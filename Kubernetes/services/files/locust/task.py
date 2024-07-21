from locust import HttpUser, task, between
import random

API_VERSION = "v1"

class MetricsUser(HttpUser):
    wait_time = between(1, 5)

    @task(1)
    def get_metrics(self):
        self.client.get("/metrics")

    @task(2)
    def list_persons(self):
        limit = random.randint(2, 5)
        params = {
            "offset": 0,
            "limit": limit
        }
        self.client.get(f"/{API_VERSION}/person", params=params)