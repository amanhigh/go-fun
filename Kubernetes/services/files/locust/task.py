from locust import HttpUser, task, between

# https://docs.locust.io/en/stable/writing-a-locustfile.html

class MetricsUser(HttpUser):
    wait_time = between(1, 5)

    @task
    def get_metrics(self):
        self.client.get("/metrics")