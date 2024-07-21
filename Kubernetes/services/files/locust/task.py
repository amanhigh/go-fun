from locust import HttpUser, task, between
import random
import json

API_VERSION = "v1"

# https://docs.locust.io/en/stable/writing-a-locustfile.html
class FunAppUser(HttpUser):
    # Wait time between tasks in seconds
    wait_time = between(1, 5)
    created_person_ids = []

    # Task with Weight 1 which is Least Frequent, Higher the number, more frequent the task
    @task(1)
    def get_metrics(self):
        self.client.get("/metrics")

    @task(2)
    def create_person(self):
        payload = {
            "name": f"User {random.randint(1, 1000)}",
            "age": random.randint(18, 80),
            "gender": random.choice(["MALE", "FEMALE"])
        }
        response = self.client.post(f"/{API_VERSION}/person", json=payload)
        if response.status_code == 201:  # Assuming 201 Created status code
            person_id = response.json().get('id')
            if person_id:
                self.created_person_ids.append(person_id)

    @task(3)
    def get_person(self):
        if self.created_person_ids:
            person_id = random.choice(self.created_person_ids)
            self.client.get(f"/{API_VERSION}/person/{person_id}")
        else:
            # If no persons have been created yet, create one
            self.create_person()
            
    @task(4)
    def list_persons(self):
        limit = random.randint(2, 5)
        params = {
            "offset": 0,
            "limit": limit
        }
        self.client.get(f"/{API_VERSION}/person", params=params)