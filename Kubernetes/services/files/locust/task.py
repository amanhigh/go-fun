from locust import HttpUser, task, between

class MyUser(HttpUser):
    wait_time = between(1, 5)

    @task
    def visit_homepage(self):
        self.client.get("/")

    @task(2)
    def visit_about_page(self):
        self.client.get("/about")

    @task(3)
    def api_request(self):
        self.client.get("/api/data")

    def on_start(self):
        # Log in at the start of each simulated user session
        self.client.post("/login", json={"username": "testuser", "password": "testpass"})
